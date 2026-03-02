package update_vacancy

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
)

type Request struct {
	CompanyID uuid.UUID
	VacancyID uuid.UUID

	Title       *string
	Description *string

	WorkFormat *value_types.WorkFormat
	City       *string

	DurationFromDays *int
	DurationToDays   *int

	EmploymentType   *value_types.EmploymentType
	HoursPerWeekFrom *int
	HoursPerWeekTo   *int

	FlexibleSchedule *bool

	IsPaid     *bool
	SalaryFrom *int
	SalaryTo   *int

	InternshipToOffer *bool
}
