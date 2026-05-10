package outbox

import (
	"context"
	"fmt"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/config"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain/events"
)

type EventType string

const (
	EventTypeResumeUpdated     EventType = "ResumeUpdated"
	EventTypeCandidateUpserted EventType = "CandidateUpserted"
)

type WriterRepository interface {
	Create(ctx context.Context, msg Message) error
}

type Writer struct {
	repo    WriterRepository
	encoder Encoder
	cfg     config.Outbox
}

func NewWriter(config config.Outbox, repository WriterRepository, encoder Encoder) *Writer {
	return &Writer{repo: repository, encoder: encoder, cfg: config}
}

func (w *Writer) WriteResumeUpdated(ctx context.Context, ev domain.ResumeUpdatedEvent) error {
	payload, schemaID, err := w.encoder.ResumeUpdatedToBytes(ev)
	if err != nil {
		return err
	}
	key := ev.Resume.ID[:]
	newMessage := Message{
		ID:            ev.EventID,
		AggregateID:   ev.Resume.ID,
		Topic:         w.cfg.ResumeTopic,
		EventType:     EventTypeResumeUpdated,
		SchemaID:      schemaID,
		Headers:       map[string]string{},
		Key:           key,
		Payload:       payload,
		Status:        StatusPending,
		AttemptCount:  0,
		CreatedAt:     ev.OccurredAt,
		NextAttemptAt: ev.OccurredAt,
	}

	err = w.repo.Create(ctx, newMessage)
	if err != nil {
		return fmt.Errorf("write vacancy published outbox: %w ", err)
	}
	return nil
}

func (w *Writer) WriteCandidateUpserted(ctx context.Context, ev events.CandidateUpserted) error {
	payload, schemaID, err := w.encoder.CandidateUpsertedToBytes(ev)
	if err != nil {
		return err
	}
	key := ev.CandidateID[:]
	newMessage := Message{
		ID:            ev.EventID,
		AggregateID:   ev.CandidateID,
		Topic:         w.cfg.CandidateTopic,
		EventType:     EventTypeCandidateUpserted,
		SchemaID:      schemaID,
		Headers:       map[string]string{},
		Key:           key,
		Payload:       payload,
		Status:        StatusPending,
		AttemptCount:  0,
		CreatedAt:     ev.OccurredAt,
		NextAttemptAt: ev.OccurredAt,
	}

	err = w.repo.Create(ctx, newMessage)
	if err != nil {
		return fmt.Errorf("write candidate upserted outbox: %w ", err)
	}
	return nil
}
