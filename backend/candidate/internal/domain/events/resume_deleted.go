package events

import (
	"time"

	"github.com/google/uuid"
)

type ResumeDeleted struct {
	EventID    uuid.UUID `avro:"event_id"`
	ResumeID   uuid.UUID `avro:"resume_id"`
	OccurredAt time.Time `avro:"occurred_at"`
}

func NewResumeDeleted(resumeId uuid.UUID) ResumeDeleted {
	return ResumeDeleted{
		EventID:    uuid.New(),
		ResumeID:   resumeId,
		OccurredAt: time.Now(),
	}
}
