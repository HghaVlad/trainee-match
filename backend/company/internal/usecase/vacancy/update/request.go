package update

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

type Request struct {
	CompanyID uuid.UUID
	VacancyID uuid.UUID

	Title       *string
	Description *string

	WorkFormat *vacancy.WorkFormat
	City       *string

	DurationFromDays *int
	DurationToDays   *int

	EmploymentType   *vacancy.EmploymentType
	HoursPerWeekFrom *int
	HoursPerWeekTo   *int

	FlexibleSchedule *bool

	IsPaid     *bool
	SalaryFrom *int
	SalaryTo   *int

	InternshipToOffer *bool
}

// lightValidation to avoid vacancy fetching,
// the main validation is the domain vacancy.Vacancy.Validate
// because of state
//
//nolint:gocognit //just validations
func (r *Request) lightValidate() error {
	// --- Title ---
	if r.Title != nil {
		l := len([]rune(*r.Title))
		if l == 0 || l > vacancy.MaxTitleLen {
			return vacancy.ErrInvalidTitleLength
		}
	}

	// --- Description ---
	if r.Description != nil {
		if len([]rune(*r.Description)) > vacancy.MaxDescriptionLen {
			return vacancy.ErrInvalidDescriptionLength
		}
	}

	// --- WorkFormat ---
	if r.WorkFormat != nil && !r.WorkFormat.IsValid() {
		return vacancy.ErrInvalidWorkFormat
	}

	// --- EmploymentType ---
	if r.EmploymentType != nil && !r.EmploymentType.IsValid() {
		return vacancy.ErrInvalidEmploymentType
	}

	// --- Salary ---
	if r.SalaryFrom != nil && *r.SalaryFrom < 0 ||
		r.SalaryTo != nil && *r.SalaryTo < 0 {
		return vacancy.ErrNegativeSalary
	}

	if r.SalaryTo != nil && *r.SalaryTo > vacancy.MaxSalary {
		return vacancy.ErrSalaryTooLarge
	}

	// --- Duration ---
	if r.DurationFromDays != nil && *r.DurationFromDays <= 0 {
		return vacancy.ErrInvalidDurationRange
	}
	if r.DurationToDays != nil && *r.DurationToDays > vacancy.MaxDurationDays {
		return vacancy.ErrInvalidDurationRange
	}

	// --- Hours ---
	if r.HoursPerWeekFrom != nil && *r.HoursPerWeekFrom <= 0 {
		return vacancy.ErrInvalidHoursRange
	}
	if r.HoursPerWeekTo != nil && *r.HoursPerWeekTo > vacancy.MaxHoursPerWeek {
		return vacancy.ErrInvalidHoursRange
	}

	return nil
}
