package eventhandler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/msgbroker/schemaregistry"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/outbox"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/projection/userhr"
)

// Handler contains the logic of handling a single event.
// Decodes it, calls inner application logic.
// If it fails, decides whether to retry, or to push it to dlq.
type Handler struct {
	cfg       config.KafkaHandling
	decoder   Decoder
	dlqSender DLQSender

	userHrCreator UserHrCreator
	searchIndexer SearchIndexer

	logger *slog.Logger
}

func NewHandler(
	cfg config.KafkaHandling,
	decoder Decoder,
	dlqSender DLQSender,
	userHrCreator UserHrCreator,
	searchIndexer SearchIndexer,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		cfg:           cfg,
		decoder:       decoder,
		dlqSender:     dlqSender,
		searchIndexer: searchIndexer,
		userHrCreator: userHrCreator,
		logger:        logger,
	}
}

// HandleMsg gets metadata common for all events, calls specific event handling.
// If unsuccessful retries, writes to dlq outbox with data that ws managed to be retrieved.
// Logs key points
func (h *Handler) HandleMsg(ctx context.Context, event *Event) {
	evType := string(event.Headers["event_type"])
	evID, err := getEventID(event.Headers)
	if err != nil {
		// maybe we will get it from payload
		h.logger.WarnContext(ctx, "event ID in headers err", "error", err)
	}

	var lastErr error

	for i := range h.cfg.MaxRetries {
		if i > 0 {
			time.Sleep(h.backoff(i)) // 50ms 100ms
		}

		res, err := h.handleByEventType(ctx, event.Payload, evType)

		switch res {
		case ResultSuccess:
			return
		case ResultDLQ:
			h.toDLQ(ctx, event, evID, evType, err.Error())
			return
		case ResultRetry:
			lastErr = err
			continue
		}
	}

	h.toDLQ(ctx, event, evID, evType,
		fmt.Sprintf("max retries: %d, lastErr: %v", h.cfg.MaxRetries, lastErr))
}

func (h *Handler) handleByEventType(ctx context.Context, payload []byte, evType string) (ResultStatus, error) {
	switch evType {
	case UserCreatedEventType:
		return h.handleUserCreated(ctx, payload)
	case string(outbox.EventTypeVacancyPublished):
		return h.handleVacancyPublished(ctx, payload)
	case string(outbox.EventTypeVacancyDraftCreated):
		return h.handleVacancyDraftCreated(ctx, payload)
	case string(outbox.EventTypeVacancyArchived):
		return h.handleVacancyArchived(ctx, payload)
	case string(outbox.EventTypeVacancyUpdated):
		return h.handleVacancyUpdated(ctx, payload)
	case string(outbox.EventTypeVacancyModerationUpdated):
		return h.handleVacancyModUpd(ctx, payload)
	case string(outbox.EventTypeCompanyUpdated):
		return h.handleCompanyUpdated(ctx, payload)
	case string(outbox.EventTypeCompanyDeleted):
		return h.handleCompanyRemoved(ctx, payload)
	default:
		return ResultDLQ, ErrUnknownEventType
	}
}

// decodes events, calls , checks errors to decide if we retry or dlq
func (h *Handler) handleUserCreated(ctx context.Context, payload []byte) (ResultStatus, error) {
	event, err := h.decoder.GetUserCreatedEvent(ctx, payload)
	if err != nil {
		return classifyErr(err), err
	}

	// filter, only handle when user is hr
	if event.Role != identity.RoleHR {
		return ResultSuccess, nil
	}

	h.logger.InfoContext(ctx, "got user created event", "id", event.UserID)

	err = h.userHrCreator.Execute(ctx, *event)

	if err != nil {
		switch {
		case errors.Is(err, userhr.ErrUserIDNil),
			errors.Is(err, userhr.ErrUsernameEmpty),
			errors.Is(err, userhr.ErrEmailEmpty):
			return ResultDLQ, err
		default:
			return ResultRetry, err
		}
	}

	return ResultSuccess, nil
}

func (h *Handler) handleVacancyPublished(ctx context.Context, payload []byte) (ResultStatus, error) {
	event, err := h.decoder.GetVacancyPublishedEvent(ctx, payload)
	if err != nil {
		return classifyErr(err), err
	}

	// TODO: maybe remove these logs
	h.logger.InfoContext(ctx, "got user vacancy pub event", "id", event.VacancyID)

	err = h.searchIndexer.Index(ctx, event.VacancyID)
	if err != nil {
		return ResultRetry, err
	}

	return ResultSuccess, nil
}

func (h *Handler) handleVacancyDraftCreated(ctx context.Context, payload []byte) (ResultStatus, error) {
	event, err := h.decoder.GetVacancyDraftCreatedEvent(ctx, payload)
	if err != nil {
		return classifyErr(err), err
	}

	err = h.searchIndexer.Index(ctx, event.VacancyID)
	if err != nil {
		return ResultRetry, err
	}

	return ResultSuccess, nil
}

func (h *Handler) handleVacancyArchived(ctx context.Context, payload []byte) (ResultStatus, error) {
	event, err := h.decoder.GetVacancyArchivedEvent(ctx, payload)
	if err != nil {
		return classifyErr(err), err
	}

	err = h.searchIndexer.Index(ctx, event.VacancyID)
	if err != nil {
		return ResultRetry, err
	}

	return ResultSuccess, nil
}

func (h *Handler) handleVacancyUpdated(ctx context.Context, payload []byte) (ResultStatus, error) {
	event, err := h.decoder.GetVacancyUpdatedEvent(ctx, payload)
	if err != nil {
		return classifyErr(err), err
	}

	err = h.searchIndexer.Index(ctx, event.VacancyID)
	if err != nil {
		return ResultRetry, err
	}

	return ResultSuccess, nil
}

func (h *Handler) handleVacancyModUpd(ctx context.Context, payload []byte) (ResultStatus, error) {
	event, err := h.decoder.GetVacancyModUpdEvent(ctx, payload)
	if err != nil {
		return classifyErr(err), err
	}

	err = h.searchIndexer.Index(ctx, event.VacancyID)
	if err != nil {
		return ResultRetry, err
	}

	return ResultSuccess, nil
}

func (h *Handler) handleCompanyUpdated(ctx context.Context, payload []byte) (ResultStatus, error) {
	event, err := h.decoder.GetCompanyUpdatedEvent(ctx, payload)
	if err != nil {
		return classifyErr(err), err
	}

	err = h.searchIndexer.UpdateCompanyName(ctx, event.CompanyID, event.CompanyName)
	if err != nil {
		return ResultRetry, err
	}

	return ResultSuccess, nil
}

func (h *Handler) handleCompanyRemoved(ctx context.Context, payload []byte) (ResultStatus, error) {
	event, err := h.decoder.GetCompanyDeletedEvent(ctx, payload)
	if err != nil {
		return classifyErr(err), err
	}

	err = h.searchIndexer.RemoveByCompany(ctx, event.CompanyID)
	if err != nil {
		return ResultRetry, err
	}

	return ResultSuccess, nil
}

func (h *Handler) toDLQ(ctx context.Context, event *Event, evID uuid.UUID, evType string, errMsg string) {
	h.logger.WarnContext(ctx, "sending to DLQ", "event_type", evType, "error", errMsg)

	var lastErr error

	for i := range 3 {
		if i > 0 {
			time.Sleep(h.backoff(i))
		}

		if err := h.dlqSender.ToDLQ(ctx, evID, event.Key, event.Payload, event.Topic, evType, errMsg); err != nil {
			lastErr = err
			continue
		}

		return
	}

	h.logger.ErrorContext(ctx, "dlq failed after retries", "error", lastErr)
}

func classifyErr(err error) ResultStatus {
	switch {
	case errors.Is(err, schemaregistry.ErrDecodePayload),
		errors.Is(err, schemaregistry.ErrSchemaNotFound),
		errors.Is(err, ErrUnknownEventType):
		return ResultDLQ
	case errors.Is(err, schemaregistry.ErrSchemaRegistryUnavailable):
		return ResultRetry
	}

	return ResultRetry
}

func (h *Handler) backoff(i int) time.Duration {
	return time.Duration(i) * h.cfg.BaseRetryDelay // 0ms 50ms 100ms
}

func getEventID(m map[string][]byte) (uuid.UUID, error) {
	b, ok := m["event_id"]
	if !ok {
		return uuid.Nil, errors.New("missing event_id")
	}

	id, err := uuid.FromBytes(b)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing event_id: %w", err)
	}

	return id, nil
}
