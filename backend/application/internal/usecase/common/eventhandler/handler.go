package eventhandler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/config"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/dlq"
)

type Handler struct {
	decoder           Decoder
	cfg               config.KafkaHandling
	dLQSender         DLQSender
	resumeUpserted    ResumeUpsertedUsecase
	ResumeDeleted     ResumeDeletedUsecase
	candidateUpserted CandidateUpsertedUsecase
}

func NewHandler(decoder Decoder, cfg config.KafkaHandling, sender DLQSender,
	resumeUpsertedUsecase ResumeUpsertedUsecase,
	resumeDeletedUsecase ResumeDeletedUsecase,
	candidateUpsertedUsecase CandidateUpsertedUsecase,
) *Handler {
	return &Handler{
		decoder:           decoder,
		cfg:               cfg,
		dLQSender:         sender,
		resumeUpserted:    resumeUpsertedUsecase,
		ResumeDeleted:     resumeDeletedUsecase,
		candidateUpserted: candidateUpsertedUsecase,
	}
}

func (h *Handler) HandleEvent(ctx context.Context, event Event) {
	eventType, ok := event.Headers["event_type"]
	if !ok {
		return
	}
	eventIDBytes, ok := event.Headers["event_id"]
	if !ok {
		return
	}
	eventID, err := uuid.FromBytes(eventIDBytes)
	if err != nil {
		return
	}

	for i := range h.cfg.RetryCount {
		if i > 0 {
			time.Sleep(h.cfg.RetryDelay * time.Duration(i)) // delay for retrying process
		}
		var status ResultStatus
		switch string(eventType) {
		case "ResumeUpserted":
			status, err = h.handleResumeUpsertedEvent(ctx, event.Payload)
		case "ResumeDeleted":
			status, err = h.handleResumeDeletedEvent(ctx, event.Payload)
		case "CandidateUpserted":
			status, err = h.handleCandidateUpsertedEvent(ctx, event.Payload)
		default:
			slog.Warn("Unknown event type: %s", eventType)
			return
		}
		if status == ResultStatusSuccess {
			return
		} else if status == ResultStatusDLQ {
			break
		}
	}
	if err == nil {
		slog.Warn("Failed to handle event with nil error but non-success status: %s", eventType)
		return
	}

	sendingErr := h.sendToDLQ(ctx, event, eventID, string(eventType), err)
	if sendingErr != nil {
		slog.Error("Failed to send to DLQ", sendingErr)
	}
}

func (h *Handler) handleResumeUpsertedEvent(ctx context.Context, payload []byte) (ResultStatus, error) {
	event, err := h.decoder.DecodeResumeUpsertedEvent(ctx, payload)
	if err != nil {
		return ResultStatusDLQ, err
	}

	err = h.resumeUpserted.Execute(ctx, event)
	if err != nil {
		return ResultStatusRetry, err
	}
	return ResultStatusSuccess, nil
}

func (h *Handler) handleResumeDeletedEvent(ctx context.Context, payload []byte) (ResultStatus, error) {
	event, err := h.decoder.DecodeResumeDeletedEvent(ctx, payload)
	if err != nil {
		return ResultStatusDLQ, err
	}

	err = h.ResumeDeleted.Execute(ctx, event)
	if err != nil {
		return ResultStatusRetry, err
	}
	return ResultStatusSuccess, nil
}

func (h *Handler) handleCandidateUpsertedEvent(ctx context.Context, payload []byte) (ResultStatus, error) {
	event, err := h.decoder.DecodeCandidateUpsertedEvent(ctx, payload)
	if err != nil {
		return ResultStatusDLQ, err
	}

	err = h.candidateUpserted.Execute(ctx, event)
	if err != nil {
		return ResultStatusRetry, err
	}
	return ResultStatusSuccess, nil
}

func (h *Handler) sendToDLQ(
	ctx context.Context,
	event Event,
	eventID uuid.UUID,
	eventType string,
	originalErr error,
) error {
	message := dlq.Message{
		EventID:           eventID,
		Payload:           event.Payload,
		OriginalTopic:     event.Topic,
		OriginalEventType: eventType,
		LastError:         originalErr.Error(),
		FailedAt:          time.Now().UTC(),
	}
	err := h.dLQSender.SendDLQ(ctx, message, event.Key)
	if err != nil {
		return fmt.Errorf("failed to send message to DLQ: %w", err)
	}
	return nil
}
