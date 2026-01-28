package create_vacancy

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
)

type Request struct {
	CompanyID uuid.UUID

	Title       string
	Description string

	WorkFormat domain.WorkFormat
	City       *string

	DurationFromMonths *int
	DurationToMonths   *int

	EmploymentType   domain.EmploymentType
	HoursPerWeekFrom *int
	HoursPerWeekTo   *int

	FlexibleSchedule bool

	IsPaid     bool
	SalaryFrom *int
	SalaryTo   *int

	InternshipToOffer bool
}
