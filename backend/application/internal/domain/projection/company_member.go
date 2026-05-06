package projection

import (
	"github.com/google/uuid"
)

type CompanyMember struct {
	UserID    uuid.UUID
	CompanyID uuid.UUID
	Role      string
}
