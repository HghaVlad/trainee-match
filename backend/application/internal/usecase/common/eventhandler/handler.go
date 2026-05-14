package eventhandler

import (
	"context"
	"log/slog"
	"time"

	"github.com/HghaVlad/trainee-match/backend/application/internal/config"
)

type Handler struct {
	decoder           Decoder
	cfg               config.KafkaHandling
	resumeUpserted    ResumeUpsertedUsecase
	ResumeDeleted     ResumeDeletedUsecase
	candidateUpserted CandidateUpsertedUsecase
}

func NewHandler(decoder Decoder, cfg config.KafkaHandling,
	resumeUpsertedUsecase ResumeUpsertedUsecase,
	resumeDeletedUsecase ResumeDeletedUsecase,
	candidateUpsertedUsecase CandidateUpsertedUsecase,
) *Handler {
	return &Handler{
		decoder:           decoder,
		cfg:               cfg,
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

	for i := range h.cfg.RetryCount {
		if i > 0 {
			time.Sleep(h.cfg.RetryDelay * time.Duration(i)) // delay for retrying process
		}
		var status ResultStatus

		switch string(eventType) {
		case "ResumeUpserted":
			status = h.handleResumeUpsertedEvent(ctx, event.Payload)
		case "ResumeDeleted":
			status = h.handleResumeDeletedEvent(ctx, event.Payload)
		case "CandidateUpserted":
			status = h.handleCandidateUpsertedEvent(ctx, event.Payload)
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

	slog.Error("Event processing failed, sending to DLQ", "event_type", string(eventType))
}

func (h *Handler) handleResumeUpsertedEvent(ctx context.Context, payload []byte) ResultStatus {
	event, err := h.decoder.DecodeResumeUpsertedEvent(ctx, payload)
	if err != nil {
		return ResultStatusDLQ
	}

	err = h.resumeUpserted.Execute(ctx, event)
	if err != nil {
		return ResultStatusDLQ
	}
	return ResultStatusSuccess
}

func (h *Handler) handleResumeDeletedEvent(ctx context.Context, payload []byte) ResultStatus {
	event, err := h.decoder.DecodeResumeDeletedEvent(ctx, payload)
	if err != nil {
		return ResultStatusDLQ
	}

	err = h.ResumeDeleted.Execute(ctx, event)
	if err != nil {
		return ResultStatusDLQ
	}
	return ResultStatusSuccess
}

func (h *Handler) handleCandidateUpsertedEvent(ctx context.Context, payload []byte) ResultStatus {
	event, err := h.decoder.DecodeCandidateUpsertedEvent(ctx, payload)
	if err != nil {
		return ResultStatusDLQ
	}

	err = h.candidateUpserted.Execute(ctx, event)
	if err != nil {
		return ResultStatusDLQ
	}
	return ResultStatusSuccess
}
