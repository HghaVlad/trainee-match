package moderationstatus

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
)

type Request struct {
	ID     uuid.UUID
	Status company.ModerationStatus
}

func (r *Request) Validate() error {
	if !r.Status.IsValid() {
		return company.ErrInvalidModerationStatus
	}

	return nil
}
