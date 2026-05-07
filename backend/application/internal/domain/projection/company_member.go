package projection

import (
	"errors"

	"github.com/google/uuid"
)

type CompanyMember struct {
	UserID    uuid.UUID
	CompanyID uuid.UUID
	Role      string
}

var (
	ErrCompanyMemberNotFound = errors.New("company member not found")
)
