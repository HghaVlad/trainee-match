package views

import (
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

type MemberVacSummary struct {
	ID             uuid.UUID
	Title          string
	WorkFormat     vacancy.WorkFormat
	City           *string
	EmploymentType vacancy.EmploymentType
	IsPaid         bool
	SalaryFrom     *int
	SalaryTo       *int
	Status         vacancy.Status
	ModStatus      vacancy.ModerationStatus
	CreatedAt      time.Time
}
