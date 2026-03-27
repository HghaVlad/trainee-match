package add

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
)

type Request struct {
	CompanyID uuid.UUID
	UserID    uuid.UUID
	Role      member.CompanyRole
}

func (r *Request) Validate() error {
	if r.UserID == uuid.Nil {
		return member.ErrInvalidUserID
	}

	if !r.Role.IsValid() {
		return member.ErrInvalidCompanyMemberRole
	}

	return nil
}
