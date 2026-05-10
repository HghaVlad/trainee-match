package events

import (
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

type ResumeDeleted struct {
	EventID    uuid.UUID `avro:"event_id"`
	ResumeID   uuid.UUID `avro:"candidate_id"`
	OccurredAt time.Time `avro:"occurred_at"`
}

func NewResumeDeleted(resume domain.Resume) ResumeDeleted {
	return ResumeDeleted{
		EventID:    uuid.New(),
		ResumeID:   resume.ID,
		OccurredAt: time.Now(),
	}
}
