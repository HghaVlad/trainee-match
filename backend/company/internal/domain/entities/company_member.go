package domain

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
)

type CompanyMember struct {
	UserID    uuid.UUID               `db:"user_id"`
	CompanyID uuid.UUID               `db:"company_id"`
	Role      value_types.CompanyRole `db:"role"`
}
