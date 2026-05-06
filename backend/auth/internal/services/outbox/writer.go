package outbox

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/config"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/domain"
)

type EventType string

const (
	EventTypeUserCreated EventType = "UserCreated"
	defaultMaxAttempts             = 5
)

type WriterRepo interface {
	Create(ctx context.Context, msg Message) error
}

type Encoder interface {
	UserCreatedToBytes(event domain.UserCreatedEvent) ([]byte, error)
}

type Writer struct {
	repo    WriterRepo
	encoder Encoder
	cfg     config.Outbox
}

func NewWriter(cfg config.Outbox, repo WriterRepo, encoder Encoder) *Writer {
	return &Writer{repo: repo, encoder: encoder, cfg: cfg}
}

func (w *Writer) WriteUserCreated(ctx context.Context, ev domain.UserCreatedEvent) error {
	payload, err := w.encoder.UserCreatedToBytes(ev)
	if err != nil {
		return fmt.Errorf("write user created outbox: %w", err)
	}

	key := ev.UserID[:]
	msg := w.createDefaultMsg(
		ev.UserID,
		payload,
		key,
		w.cfg.UserTopic,
		EventTypeUserCreated,
		ev.EventID,
		ev.OccurredAt,
	)

	if err := w.repo.Create(ctx, msg); err != nil {
		return fmt.Errorf("write user created outbox: %w", err)
	}

	return nil
}

func (w *Writer) createDefaultMsg(
	aggregateID uuid.UUID,
	payload, key []byte,
	topic string,
	evType EventType,
	eventID uuid.UUID,
	occurredAt time.Time,
) Message {
	return Message{
		ID:            eventID,
		AggregateID:   aggregateID,
		Topic:         topic,
		Key:           key,
		Payload:       payload,
		Headers:       make(map[string]string),
		SchemaID:      schemaIDFromPayload(payload),
		EventType:     evType,
		Status:        StatusPending,
		MaxAttempts:   defaultMaxAttempts,
		CreatedAt:     occurredAt,
		NextAttemptAt: occurredAt,
	}
}

func schemaIDFromPayload(payload []byte) int {
	return int(binary.BigEndian.Uint32(payload[1:]))
}
