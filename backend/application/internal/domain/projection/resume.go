package projection

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Resume struct {
	ID          uuid.UUID
	CandidateID uuid.UUID
	Name        string
	Data        ResumeData
	Status      ResumeStatus
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

var (
	ErrResumeNotFound = errors.New("resume projection not found")
)

type ResumeUpsertedEvent struct {
	EventID     uuid.UUID    `avro:"event_id"`
	ResumeID    uuid.UUID    `avro:"resume_id"`
	OccurredAt  time.Time    `avro:"occurred_at"`
	CandidateID uuid.UUID    `avro:"candidate_id"`
	Name        string       `avro:"name"`
	Data        ResumeData   `avro:"data"`
	Status      ResumeStatus `avro:"status"`
	CreatedAt   *time.Time   `avro:"created_at"`
	UpdatedAt   *time.Time   `avro:"updated_at"`
}

func (ev ResumeUpsertedEvent) ToResume() Resume {
	return Resume{
		ID:          ev.ResumeID,
		CandidateID: ev.CandidateID,
		Name:        ev.Name,
		Data:        ev.Data,
		Status:      ev.Status,
		CreatedAt:   ev.CreatedAt,
		UpdatedAt:   ev.UpdatedAt,
	}
}

type ResumeDeletedEvent struct {
	EventID    uuid.UUID `avro:"event_id"`
	ResumeID   uuid.UUID `avro:"resume_id"`
	OccurredAt time.Time `avro:"occurred_at"`
}
