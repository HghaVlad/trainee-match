package member

import (
	"github.com/google/uuid"
)

type CompanyMember struct {
	UserID    uuid.UUID   `db:"user_id"`
	CompanyID uuid.UUID   `db:"company_id"`
	Role      CompanyRole `db:"role"`
}
