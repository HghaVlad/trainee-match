package list

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
)

type View struct {
	UserID    uuid.UUID
	CompanyID uuid.UUID
	Username  string
	Email     string
	Role      member.CompanyRole
}

type Response struct {
	Members []View
	HasMore bool
}
