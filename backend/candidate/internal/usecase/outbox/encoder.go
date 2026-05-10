package outbox

import (
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain/events"
)

type Encoder interface {
	ResumeUpdatedToBytes(ev domain.ResumeUpdatedEvent) ([]byte, int, error)
	CandidateUpsertedToBytes(ev events.CandidateUpserted) ([]byte, int, error)
}
