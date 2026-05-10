package events

import (
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

type CandidateDeleted struct {
	EventID     uuid.UUID `avro:"event_id"`
	CandidateID uuid.UUID `avro:"candidate_id"`
	OccurredAt  time.Time `avro:"occurred_at"`
}

func NewCandidateDeleted(c domain.Candidate) CandidateDeleted {
	evt := CandidateDeleted{
		EventID:     uuid.New(),
		CandidateID: c.ID,
		OccurredAt:  time.Now().UTC(),
	}
	return evt
}
