package outbox

import "github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"

type Encoder interface {
	ResumeUpdatedToBytes(ev domain.ResumeUpdatedEvent) ([]byte, int, error)
}
