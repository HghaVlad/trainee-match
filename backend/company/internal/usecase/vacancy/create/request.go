package create

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

type Request struct {
	CompanyID uuid.UUID

	Title       string
	Description string

	WorkFormat vacancy.WorkFormat
	City       *string

	DurationFromDays *int
	DurationToDays   *int

	EmploymentType   *vacancy.EmploymentType
	HoursPerWeekFrom *int
	HoursPerWeekTo   *int

	FlexibleSchedule bool

	IsPaid     bool
	SalaryFrom *int
	SalaryTo   *int

	InternshipToOffer bool
}
