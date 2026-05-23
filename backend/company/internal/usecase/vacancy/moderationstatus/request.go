package moderationstatus

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

type Request struct {
	ID     uuid.UUID
	Status vacancy.ModerationStatus
}

func (r *Request) Validate() error {
	if !r.Status.IsValid() {
		return vacancy.ErrInvalidModerationStatus
	}

	return nil
}
