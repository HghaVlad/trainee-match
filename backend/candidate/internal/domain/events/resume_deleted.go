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

func NewResumeDeleted(resumeID uuid.UUID) ResumeDeleted {
	return ResumeDeleted{
		EventID:    uuid.New(),
		ResumeID:   resumeID,
		OccurredAt: time.Now().UTC(),
	}
}
