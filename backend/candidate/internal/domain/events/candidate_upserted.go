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

func NewCandidateUpserted(c domain.Candidate, email string) CandidateUpserted {
	evt := CandidateUpserted{
		EventID:     uuid.New(),
		CandidateID: c.ID,
		FullName:    c.FullName,
		Email:       email,
		OccurredAt:  time.Now().UTC(),
	}

	if c.Telegram != "" {
		t := c.Telegram
		evt.Telegram = &t
	}

	return evt
}
