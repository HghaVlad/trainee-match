package outbox

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

type EventType string

const (
	EventTypeVacancyPublished EventType = "VacancyPublished"
	EventTypeVacancyArchived  EventType = "VacancyArchived"
	EventTypeVacancyUpdated   EventType = "VacancyUpdated"

	EventTypeRecruiterAdded   EventType = "CompanyMemberAdded"
	EventTypeRecruiterRemoved EventType = "CompanyMemberRemoved"

	EventTypeCompanyDeleted EventType = "CompanyDeleted"
	EventTypeCompanyUpdated EventType = "CompanyUpdated"
)

const (
	VacancyTopic   = "vacancy.events"
	RecruiterTopic = "companymember.events"
	CompanyTopic   = "company.events"
)

const defaultMaxAttempts = 5

type Writer struct {
	repo    WriterRepo
	encoder Encoder
}

func NewWriter(repo WriterRepo, encoder Encoder) *Writer {
	return &Writer{
		repo:    repo,
		encoder: encoder,
	}
}

func (w *Writer) WriteVacancyPublished(ctx context.Context, ev vacancy.PublishedEvent) error {
	payload, err := w.encoder.VacancyPublishedToBytes(ev)
	if err != nil {
		return fmt.Errorf("write vacancy published outbox: %w ", err)
	}

	key := ev.VacancyID[:]
	msg := w.createDefaultMsg(
		payload,
		key,
		VacancyTopic,
		EventTypeVacancyPublished,
		ev.EventID,
		ev.OccurredAt,
	)

	err = w.repo.Create(ctx, msg)
	if err != nil {
		return fmt.Errorf("write vacancy published outbox: %w ", err)
	}
	return nil
}

func (w *Writer) WriteVacancyArchived(ctx context.Context, ev vacancy.ArchivedEvent) error {
	payload, err := w.encoder.VacancyArchivedToBytes(ev)
	if err != nil {
		return fmt.Errorf("write vacancy archived outbox: %w ", err)
	}

	key := ev.VacancyID[:]
	msg := w.createDefaultMsg(
		payload,
		key,
		VacancyTopic,
		EventTypeVacancyArchived,
		ev.EventID,
		ev.OccurredAt,
	)

	err = w.repo.Create(ctx, msg)
	if err != nil {
		return fmt.Errorf("write vacancy archived outbox: %w ", err)
	}
	return nil
}

func (w *Writer) WriteVacancyUpdated(ctx context.Context, ev vacancy.UpdatedEvent) error {
	payload, err := w.encoder.VacancyUpdatedToBytes(ev)
	if err != nil {
		return fmt.Errorf("write vacancy updated outbox: %w ", err)
	}

	key := ev.VacancyID[:]
	msg := w.createDefaultMsg(
		payload,
		key,
		VacancyTopic,
		EventTypeVacancyUpdated,
		ev.EventID,
		ev.OccurredAt,
	)

	err = w.repo.Create(ctx, msg)
	if err != nil {
		return fmt.Errorf("write vacancy updated outbox: %w ", err)
	}
	return nil
}

func (w *Writer) WriteCompanyMemberAdded(ctx context.Context, ev member.AddedEvent) error {
	payload, err := w.encoder.CompanyMemberAddedToBytes(ev)
	if err != nil {
		return fmt.Errorf("write company added outbox: %w ", err)
	}

	key := ev.CompanyID[:]
	msg := w.createDefaultMsg(
		payload,
		key,
		RecruiterTopic,
		EventTypeRecruiterAdded,
		ev.EventID,
		ev.OccurredAt,
	)

	err = w.repo.Create(ctx, msg)
	if err != nil {
		return fmt.Errorf("write company member added outbox: %w ", err)
	}
	return nil
}

func (w *Writer) WriteCompanyMemberRemoved(ctx context.Context, ev member.RemovedEvent) error {
	payload, err := w.encoder.CompanyMemberRemovedToBytes(ev)
	if err != nil {
		return fmt.Errorf("write company member removed outbox: %w ", err)
	}

	key := ev.CompanyID[:]
	msg := w.createDefaultMsg(
		payload,
		key,
		RecruiterTopic,
		EventTypeRecruiterRemoved,
		ev.EventID,
		ev.OccurredAt,
	)

	err = w.repo.Create(ctx, msg)
	if err != nil {
		return fmt.Errorf("write company member removed outbox: %w ", err)
	}
	return nil
}

func (w *Writer) WriteCompanyUpdated(ctx context.Context, ev company.UpdatedEvent) error {
	payload, err := w.encoder.CompanyUpdatedToBytes(ev)
	if err != nil {
		return fmt.Errorf("write recruiter added outbox: %w ", err)
	}

	key := ev.CompanyID[:]
	msg := w.createDefaultMsg(
		payload,
		key,
		CompanyTopic,
		EventTypeCompanyUpdated,
		ev.EventID,
		ev.OccurredAt,
	)

	err = w.repo.Create(ctx, msg)
	if err != nil {
		return fmt.Errorf("write recruiter added outbox: %w ", err)
	}
	return nil
}

func (w *Writer) WriteCompanyDeleted(ctx context.Context, ev company.DeletedEvent) error {
	payload, err := w.encoder.CompanyDeletedToBytes(ev)
	if err != nil {
		return fmt.Errorf("write company deleted outbox: %w ", err)
	}

	key := ev.CompanyID[:]
	msg := w.createDefaultMsg(
		payload,
		key,
		CompanyTopic,
		EventTypeCompanyDeleted,
		ev.EventID,
		ev.OccurredAt,
	)

	err = w.repo.Create(ctx, msg)
	if err != nil {
		return fmt.Errorf("write company deleted outbox: %w ", err)
	}
	return nil
}

func (w *Writer) createDefaultMsg(
	payload, key []byte,
	topic string,
	evType EventType,
	id uuid.UUID,
	occurredAt time.Time,
) Message {
	return Message{
		ID:            id,
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
