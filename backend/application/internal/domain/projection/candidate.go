package projection

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Candidate struct {
	ID        uuid.UUID
	FullName  string
	Email     string
	Telegram  *string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

var (
	ErrCandidateNotFound = errors.New("candidate projection not found")
)

type CandidateUpsertedEvent struct {
	EventID     uuid.UUID `avro:"event_id"`
	CandidateID uuid.UUID `avro:"candidate_id"`
	FullName    string    `avro:"full_name"`
	Email       string    `avro:"email"`
	Telegram    *string   `avro:"telegram"`
	OccurredAt  time.Time `avro:"occurred_at"`
}

func (c CandidateUpsertedEvent) ToCandidate() Candidate {
	return Candidate{
		ID:       c.CandidateID,
		FullName: c.FullName,
		Email:    c.Email,
		Telegram: c.Telegram,
	}
}
