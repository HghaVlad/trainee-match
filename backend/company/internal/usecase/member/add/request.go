package add

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/projection/userhr"
)

type Request struct {
	CompanyID uuid.UUID
	Username  string
	Role      member.CompanyRole
}

func (r *Request) Validate() error {
	if r.Username == "" {
		return userhr.ErrUsernameEmpty
	}

	if !r.Role.IsValid() {
		return member.ErrInvalidCompanyMemberRole
	}

	return nil
}
