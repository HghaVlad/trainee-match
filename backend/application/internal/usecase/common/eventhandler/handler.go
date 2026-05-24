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
	decoder              Decoder
	cfg                  config.KafkaHandling
	logger               *slog.Logger
	dLQSender            DLQSender
	resumeUpserted       ResumeUpsertedUsecase
	resumeDeleted        ResumeDeletedUsecase
	resumeArchived       ResumeArchivedUsecase
	candidateUpserted    CandidateUpsertedUsecase
	companyUpdated       CompanyUpdatedUsecase
	companyDeleted       CompanyDeletedUsecase
	companyMemberAdded   CompanyMemberAddedUsecase
	companyMemberRemoved CompanyMemberRemovedUsecase
	vacancyPublished     VacancyPublishedUsecase
	vacancyArchived      VacancyArchivedUsecase
	vacancyUpdated       VacancyUpdatedUsecase
	vacancyModUpd        VacancyModerationUpdUsecase
	companyModUpd        CompanyModerationUpdUsecase
}

func NewHandler(
	decoder Decoder,
	cfg config.KafkaHandling,
	logger *slog.Logger,
	sender DLQSender,
	resumeUpsertedUsecase ResumeUpsertedUsecase,
	resumeDeletedUsecase ResumeDeletedUsecase,
	resumeArchivedUsecase ResumeArchivedUsecase,
	candidateUpsertedUsecase CandidateUpsertedUsecase,
	companyUpdatedUsecase CompanyUpdatedUsecase,
	companyDeletedUsecase CompanyDeletedUsecase,
	companyMemberAddedUsecase CompanyMemberAddedUsecase,
	companyMemberRemovedUsecase CompanyMemberRemovedUsecase,
	vacancyPublishedUsecase VacancyPublishedUsecase,
	vacancyArchivedUsecase VacancyArchivedUsecase,
	vacancyUpdatedUsecase VacancyUpdatedUsecase,
	vacancyModUpd VacancyModerationUpdUsecase,
	companyModUpd CompanyModerationUpdUsecase,
) *Handler {
	return &Handler{
		decoder:              decoder,
		cfg:                  cfg,
		logger:               logger,
		dLQSender:            sender,
		resumeUpserted:       resumeUpsertedUsecase,
		resumeDeleted:        resumeDeletedUsecase,
		resumeArchived:       resumeArchivedUsecase,
		candidateUpserted:    candidateUpsertedUsecase,
		companyUpdated:       companyUpdatedUsecase,
		companyDeleted:       companyDeletedUsecase,
		companyMemberAdded:   companyMemberAddedUsecase,
		companyMemberRemoved: companyMemberRemovedUsecase,
		vacancyPublished:     vacancyPublishedUsecase,
		vacancyArchived:      vacancyArchivedUsecase,
		vacancyUpdated:       vacancyUpdatedUsecase,
		vacancyModUpd:        vacancyModUpd,
		companyModUpd:        companyModUpd,
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

	for i := range min(h.cfg.RetryCount, 1) { // doing at least once
		if i > 0 {
			time.Sleep(h.cfg.RetryDelay * time.Duration(i)) // delay for retrying process
		}
		var status ResultStatus
		switch string(eventType) {
		case "ResumeUpserted":
			status, err = h.handleResumeUpsertedEvent(ctx, event.Payload)
		case "ResumeDeleted":
			status, err = h.handleResumeDeletedEvent(ctx, event.Payload)
		case "ResumeArchived":
			status, err = h.handleResumeArchivedEvent(ctx, event.Payload)
		case "CandidateUpserted":
			status, err = h.handleCandidateUpsertedEvent(ctx, event.Payload)
		case "CompanyUpdated":
			status, err = h.handleCompanyUpdatedEvent(ctx, event.Payload)
		case "CompanyDeleted":
			status, err = h.handleCompanyDeletedEvent(ctx, event.Payload)
		case "CompanyMemberAdded":
			status, err = h.handleCompanyMemberAddedEvent(ctx, event.Payload)
		case "CompanyMemberRemoved":
			status, err = h.handleCompanyMemberRemovedEvent(ctx, event.Payload)
		case "VacancyPublished":
			status, err = h.handleVacancyPublishedEvent(ctx, event.Payload)
		case "VacancyArchived":
			status, err = h.handleVacancyArchivedEvent(ctx, event.Payload)
		case "VacancyUpdated":
			status, err = h.handleVacancyUpdatedEvent(ctx, event.Payload)
		case "VacancyModerationUpdated":
			status, err = h.handleVacModUpdEvent(ctx, event.Payload)
		case "CompanyModerationUpdated":
			status, err = h.handleCompModUpdEvent(ctx, event.Payload)
		case "VacancyDraftCreated":
			status = ResultStatusSuccess
		default:
			h.logger.WarnContext(ctx, "unknown event type", "eventType", eventType)
			return
		}
		if status == ResultStatusSuccess {
			return
		} else if status == ResultStatusDLQ {
			break
		}
	}
	if err == nil {
		h.logger.WarnContext(
			ctx,
			"failed to handle event with nil error but non-success status",
			"eventType",
			eventType,
		)
		return
	}

	sendingErr := h.sendToDLQ(ctx, event, eventID, string(eventType), err)
	if sendingErr != nil {
		h.logger.ErrorContext(ctx, "failed to send to DLQ", "err", sendingErr)
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

	err = h.resumeDeleted.Execute(ctx, event)
	if err != nil {
		return ResultStatusRetry, err
	}
	return ResultStatusSuccess, nil
}

func (h *Handler) handleResumeArchivedEvent(ctx context.Context, payload []byte) (ResultStatus, error) {
	event, err := h.decoder.DecodeResumeArchivedEvent(ctx, payload)
	if err != nil {
		return ResultStatusDLQ, err
	}

	err = h.resumeArchived.Execute(ctx, event)
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

func (h *Handler) handleCompanyUpdatedEvent(ctx context.Context, payload []byte) (ResultStatus, error) {
	event, err := h.decoder.DecodeCompanyUpdatedEvent(ctx, payload)
	if err != nil {
		return ResultStatusDLQ, err
	}

	err = h.companyUpdated.Execute(ctx, event)
	if err != nil {
		return ResultStatusRetry, err
	}
	return ResultStatusSuccess, nil
}

func (h *Handler) handleCompanyDeletedEvent(ctx context.Context, payload []byte) (ResultStatus, error) {
	event, err := h.decoder.DecodeCompanyDeletedEvent(ctx, payload)
	if err != nil {
		return ResultStatusDLQ, err
	}

	err = h.companyDeleted.Execute(ctx, event)
	if err != nil {
		return ResultStatusRetry, err
	}
	return ResultStatusSuccess, nil
}

func (h *Handler) handleCompanyMemberAddedEvent(ctx context.Context, payload []byte) (ResultStatus, error) {
	event, err := h.decoder.DecodeCompanyMemberAddedEvent(ctx, payload)
	if err != nil {
		return ResultStatusDLQ, err
	}

	err = h.companyMemberAdded.Execute(ctx, event)
	if err != nil {
		return ResultStatusRetry, err
	}
	return ResultStatusSuccess, nil
}

func (h *Handler) handleCompanyMemberRemovedEvent(ctx context.Context, payload []byte) (ResultStatus, error) {
	event, err := h.decoder.DecodeCompanyMemberRemovedEvent(ctx, payload)
	if err != nil {
		return ResultStatusDLQ, err
	}

	err = h.companyMemberRemoved.Execute(ctx, event)
	if err != nil {
		return ResultStatusRetry, err
	}
	return ResultStatusSuccess, nil
}

func (h *Handler) handleVacancyPublishedEvent(ctx context.Context, payload []byte) (ResultStatus, error) {
	event, err := h.decoder.DecodeVacancyPublishedEvent(ctx, payload)
	if err != nil {
		return ResultStatusDLQ, err
	}

	err = h.vacancyPublished.Execute(ctx, event)
	if err != nil {
		return ResultStatusRetry, err
	}
	return ResultStatusSuccess, nil
}

func (h *Handler) handleVacancyArchivedEvent(ctx context.Context, payload []byte) (ResultStatus, error) {
	event, err := h.decoder.DecodeVacancyArchivedEvent(ctx, payload)
	if err != nil {
		return ResultStatusDLQ, err
	}

	err = h.vacancyArchived.Execute(ctx, event)
	if err != nil {
		return ResultStatusRetry, err
	}
	return ResultStatusSuccess, nil
}

func (h *Handler) handleVacancyUpdatedEvent(ctx context.Context, payload []byte) (ResultStatus, error) {
	event, err := h.decoder.DecodeVacancyUpdatedEvent(ctx, payload)
	if err != nil {
		return ResultStatusDLQ, err
	}

	err = h.vacancyUpdated.Execute(ctx, event)
	if err != nil {
		return ResultStatusRetry, err
	}
	return ResultStatusSuccess, nil
}

func (h *Handler) handleVacModUpdEvent(ctx context.Context, payload []byte) (ResultStatus, error) {
	event, err := h.decoder.DecodeVacancyModerationUpdEvent(ctx, payload)
	if err != nil {
		return ResultStatusDLQ, err
	}

	err = h.vacancyModUpd.Execute(ctx, event)
	if err != nil {
		return ResultStatusRetry, err
	}

	return ResultStatusSuccess, nil
}

func (h *Handler) handleCompModUpdEvent(ctx context.Context, payload []byte) (ResultStatus, error) {
	event, err := h.decoder.DecodeCompanyModerationUpdEvent(ctx, payload)
	if err != nil {
		return ResultStatusDLQ, err
	}

	err = h.companyModUpd.Execute(ctx, event)
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
