package eventhandler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/msgbroker/schemaregistry"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"

	"github.com/HghaVlad/trainee-match/backend/company/internal/config"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/outbox"
)

// Handler contains the logic of handling a single event.
// Decodes it, calls inner application logic.
// If it fails, decides whether to retry, or to push it to dlq.
type Handler struct {
	cfg          config.KafkaHandling
	decoder      decoder
	outboxWriter outboxDLQWriter
	logger       *slog.Logger
}

func NewHandler(
	cfg config.KafkaHandling,
	decoder decoder,
	outboxWriter outboxDLQWriter,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		cfg:          cfg,
		decoder:      decoder,
		outboxWriter: outboxWriter,
		logger:       logger,
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
		h.logger.Warn("event ID in headers err", "error", err)
	}

	schemaID, err := h.decoder.RetrieveSchemaID(event.Payload)
	if err != nil {
		h.toDLQ(ctx, event, evID, schemaID, evType, event.Headers, err.Error())
		return
	}

	for i := range h.cfg.MaxRetries {
		if i > 0 {
			time.Sleep(h.backoff(i)) // 50ms 100ms
		}

		res, err := h.handleByEventType(ctx, schemaID, event.Payload, evType)

		switch res {
		case ResultSuccess:
			return
		case ResultDLQ:
			h.toDLQ(ctx, event, evID, schemaID, evType, event.Headers, err.Error())
			return
		case ResultRetry:
			continue
		}
	}

	h.logger.Warn("sending to DLQ due to max retries", "retries", h.cfg.MaxRetries)
	h.toDLQ(ctx, event, evID, schemaID, evType, event.Headers, "max retries")
}

func (h *Handler) handleByEventType(
	ctx context.Context,
	schemaID int,
	payload []byte,
	evType string,
) (ResultStatus, error) {
	switch evType {
	case UserCreatedEventType:
		return h.handleUserCreated(ctx, schemaID, payload)
	default:
		return ResultDLQ, ErrUnknownEventType
	}
}

// decodes events, calls , checks errors to decide if we retry or dlq
func (h *Handler) handleUserCreated(ctx context.Context, schemaID int, payload []byte) (ResultStatus, error) {
	event, err := h.decoder.GetUserCreatedEvent(ctx, schemaID, payload)
	if err != nil {
		return classifyErr(err), err
	}

	// filter, only handle when user is hr
	if event.Role != identity.RoleHR {
		return ResultSuccess, nil
	}

	h.logger.Info("got user created event", "id", event.EventID, "")

	// TODO: call actual usecase (better by interface)

	if err != nil {
		// return retry / dlq /success depending on error
	}

	return ResultSuccess, nil
}

func (h *Handler) toDLQ(
	ctx context.Context, msg *Event, eventID uuid.UUID,
	schemaID int, evType string,
	headers map[string][]byte, errMsg string,
) {
	h.logger.Warn("sending to DLQ", "event_type", evType, "error", errMsg)

	meta := outbox.DLQMeta{
		EventID:   eventID,
		EventType: evType,
		Topic:     msg.Topic,
		Key:       msg.Key,
		Payload:   msg.Payload,
		SchemaID:  schemaID,
		ErrMsg:    errMsg,
		Headers:   headers,
	}

	err := h.outboxWriter.WriteToDLQ(ctx, meta)
	if err != nil {
		h.logger.Warn("failed to send DLQ msg", "error", err)
	}
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
