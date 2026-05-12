package outbox

import (
	"context"
	"fmt"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/config"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain/events"
)

type EventType string

const (
	EventTypeResumeDeleted     EventType = "ResumeDeleted"
	EventTypeResumeUpserted    EventType = "ResumeUpserted"
	EventTypeCandidateUpserted EventType = "CandidateUpserted"
	EventTypeCandidateDeleted  EventType = "CandidateDeleted"
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

func (w *Writer) WriteResumeUpserted(ctx context.Context, ev events.ResumeUpserted) error {
	payload, schemaID, err := w.encoder.ResumeUpdatedToBytes(ev)
	if err != nil {
		return err
	}
	key := ev.ResumeID[:]
	newMessage := Message{
		ID:            ev.EventID,
		AggregateID:   ev.ResumeID,
		Topic:         w.cfg.ResumeTopic,
		EventType:     EventTypeResumeUpserted,
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
		return fmt.Errorf("write resume upserted outbox: %w ", err)
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

func (w *Writer) WriteResumeDeleted(ctx context.Context, ev events.ResumeDeleted) error {
	payload, schemaID, err := w.encoder.ResumeDeletedToBytes(ev)
	if err != nil {
		return err
	}
	key := ev.ResumeID[:]
	newMessage := Message{
		ID:            ev.EventID,
		AggregateID:   ev.ResumeID,
		Topic:         w.cfg.ResumeTopic,
		EventType:     EventTypeResumeDeleted,
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
		return fmt.Errorf("write resume deleted outbox: %w ", err)
	}
	return nil
}

func (w *Writer) WriteCandidateDeleted(ctx context.Context, ev events.CandidateDeleted) error {
	payload, schemaID, err := w.encoder.CandidateDeletedToBytes(ev)
	if err != nil {
		return err
	}
	key := ev.CandidateID[:]
	newMessage := Message{
		ID:            ev.EventID,
		AggregateID:   ev.CandidateID,
		Topic:         w.cfg.CandidateTopic,
		EventType:     EventTypeCandidateDeleted,
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
		return fmt.Errorf("write candidate deleted outbox: %w ", err)
	}
	return nil
}
