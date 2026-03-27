package vacancy

import (
	"time"

	"github.com/google/uuid"
)

type Vacancy struct {
	ID        uuid.UUID
	CompanyID uuid.UUID

	Title       string
	Description string

	WorkFormat WorkFormat
	City       *string

	DurationFromDays *int
	DurationToDays   *int

	EmploymentType   EmploymentType
	HoursPerWeekFrom *int
	HoursPerWeekTo   *int

	FlexibleSchedule bool

	IsPaid     bool
	SalaryFrom *int
	SalaryTo   *int

	InternshipToOffer bool

	Status      VacancyStatus
	PublishedAt *time.Time

	CreatedBy uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

const (
	maxTitleLen       = 200
	maxDescriptionLen = 5000
)

const (
	MaxSalary       = 15_000_000
	MaxDurationDays = 1800
	MaxHoursPerWeek = 80
)

// Validate checks domain invariants
func (v *Vacancy) Validate() error {
	if !v.WorkFormat.IsValid() {
		return ErrInvalidWorkFormat
	}

	if !v.EmploymentType.IsValid() {
		return ErrInvalidEmploymentType
	}

	if v.DurationFromDays != nil && v.DurationToDays != nil {
		if *v.DurationFromDays > *v.DurationToDays {
			return ErrInvalidDurationRange
		}
	}

	if v.HoursPerWeekFrom != nil && v.HoursPerWeekTo != nil {
		if *v.HoursPerWeekFrom > *v.HoursPerWeekTo {
			return ErrInvalidHoursRange
		}
	}

	if v.SalaryFrom != nil && v.SalaryTo != nil {
		if *v.SalaryFrom > *v.SalaryTo {
			return ErrInvalidSalaryRange
		}
	}

	if !v.IsPaid {
		if v.SalaryFrom != nil || v.SalaryTo != nil {
			return ErrSalaryProvidedForUnpaid
		}
	}

	if v.SalaryTo != nil && *v.SalaryTo > MaxSalary {
		return ErrSalaryTooLarge
	}

	if v.SalaryFrom != nil && *v.SalaryFrom < 0 {
		return ErrNegativeSalary
	}

	if v.DurationFromDays != nil && *v.DurationFromDays <= 0 ||
		v.DurationToDays != nil && *v.DurationToDays > MaxDurationDays {
		return ErrInvalidDurationRange
	}

	if v.HoursPerWeekFrom != nil && *v.HoursPerWeekFrom <= 0 ||
		v.HoursPerWeekTo != nil && *v.HoursPerWeekTo > MaxHoursPerWeek {
		return ErrInvalidHoursRange
	}

	if len([]rune(v.Title)) == 0 || len(v.Title) > maxTitleLen {
		return ErrInvalidTitleLength
	}

	if len([]rune(v.Description)) > maxDescriptionLen {
		return ErrInvalidDescriptionLength
	}

	return nil
}
