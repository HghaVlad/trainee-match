package events

import (
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

type CandidateUpserted struct {
	EventID     uuid.UUID `avro:"event_id"`
	CandidateID uuid.UUID `avro:"candidate_id"`
	FullName    string    `avro:"full_name"`
	Email       string    `avro:"email"`
	Telegram    *string   `avro:"telegram"`
	OccurredAt  time.Time `avro:"occurred_at"`
}

// NewCandidateUpsertedEvent создаёт событие из доменной сущности Candidate
func NewCandidateUpserted(c domain.Candidate) *CandidateUpserted {
	evt := &CandidateUpserted{
		EventID:     uuid.New(),
		CandidateID: c.ID,
		FullName:    "full name",
		Email:       "email",
		OccurredAt:  time.Now().UTC(),
	}

	if c.Telegram != "" {
		t := c.Telegram
		evt.Telegram = &t
	}

	return evt
}
