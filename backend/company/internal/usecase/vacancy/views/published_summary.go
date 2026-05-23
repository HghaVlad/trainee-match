package views

import (
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

type PublishedVacSummary struct {
	ID             uuid.UUID
	CompanyID      uuid.UUID
	CompanyName    string
	Title          string
	WorkFormat     vacancy.WorkFormat
	City           *string
	EmploymentType vacancy.EmploymentType
	IsPaid         bool
	SalaryFrom     *int
	SalaryTo       *int
	PublishedAt    time.Time
}
