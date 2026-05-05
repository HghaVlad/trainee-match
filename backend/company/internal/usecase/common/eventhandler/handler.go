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
)

// Handler contains the logic of handling a single event.
// Decodes it, calls inner application logic.
// If it fails, decides whether to retry, or to push it to dlq.
type Handler struct {
	cfg       config.KafkaHandling
	decoder   Decoder
	dlqSender DLQSender
	logger    *slog.Logger
}

func NewHandler(
	cfg config.KafkaHandling,
	decoder Decoder,
	dlqSender DLQSender,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		cfg:       cfg,
		decoder:   decoder,
		dlqSender: dlqSender,
		logger:    logger,
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

	h.logger.Info("got user created event", "id", event.UserID)

	// TODO: call actual usecase (better by interface)

	if err != nil {
		// return retry / dlq /success depending on error
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
