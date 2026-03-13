package add_member

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
)

type Request struct {
	CompanyID uuid.UUID
	UserID    uuid.UUID
	Role      value_types.CompanyRole
}

func (r *Request) Validate() error {
	if r.UserID == uuid.Nil {
		return domain_errors.ErrInvalidUserID
	}

	if !r.Role.IsValid() {
		return domain_errors.ErrInvalidCompanyMemberRole
	}

	return nil
}
