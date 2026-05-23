package outbox

import (
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain/events"
)

type Encoder interface {
	ResumeUpdatedToBytes(ev events.ResumeUpserted) ([]byte, int, error)
	CandidateUpsertedToBytes(ev events.CandidateUpserted) ([]byte, int, error)
	ResumeDeletedToBytes(ev events.ResumeDeleted) ([]byte, int, error)
	CandidateDeletedToBytes(ev events.CandidateDeleted) ([]byte, int, error)
	ResumeArchivedToBytes(ev events.ResumeArchived) ([]byte, int, error)
}
