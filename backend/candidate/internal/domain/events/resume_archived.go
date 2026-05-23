package events

import (
	"time"

	"github.com/google/uuid"
)

type ResumeArchived struct {
	EventID    uuid.UUID `avro:"event_id"`
	ResumeID   uuid.UUID `avro:"resume_id"`
	OccurredAt time.Time `avro:"occurred_at"`
}

func NewResumeArchived(resumeID uuid.UUID) ResumeArchived {
	return ResumeArchived{
		EventID:    uuid.New(),
		ResumeID:   resumeID,
		OccurredAt: time.Now().UTC(),
	}
}
