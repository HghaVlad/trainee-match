package domain

import (
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
)

type Vacancy struct {
	ID        uuid.UUID `db:"id"`
	CompanyID uuid.UUID `db:"company_id"`

	Title       string `db:"title"`
	Description string `db:"description"`

	WorkFormat value_types.WorkFormat `db:"work_format"`
	City       *string                `db:"city"`

	DurationFromMonths *int `db:"duration_from_months"`
	DurationToMonths   *int `db:"duration_to_months"`

	EmploymentType   value_types.EmploymentType `db:"employment_type"`
	HoursPerWeekFrom *int                       `db:"hours_per_week_from"`
	HoursPerWeekTo   *int                       `db:"hours_per_week_to"`

	FlexibleSchedule bool `db:"flexible_schedule"`

	IsPaid     bool `db:"is_paid"`
	SalaryFrom *int `db:"salary_from"`
	SalaryTo   *int `db:"salary_to"`

	InternshipToOffer bool `db:"internship_to_offer"`

	IsActive    bool      `db:"is_active"`
	PublishedAt time.Time `db:"published_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAtAt time.Time `db:"updated_at"`
}

const (
	maxTitleLen       = 200
	maxDescriptionLen = 5000
)

const (
	maxSalary         = 15_000_000
	maxDurationMonths = 60
	maxHoursPerWeek   = 80
)

// Validate checks domain invariants
func (v *Vacancy) Validate() error {
	if !v.WorkFormat.IsValid() {
		return domain_errors.ErrInvalidWorkFormat
	}

	if !v.EmploymentType.IsValid() {
		return domain_errors.ErrInvalidEmploymentType
	}

	if v.DurationFromMonths != nil && v.DurationToMonths != nil {
		if *v.DurationFromMonths > *v.DurationToMonths {
			return domain_errors.ErrInvalidDurationRange
		}
	}

	if v.HoursPerWeekFrom != nil && v.HoursPerWeekTo != nil {
		if *v.HoursPerWeekFrom > *v.HoursPerWeekTo {
			return domain_errors.ErrInvalidHoursRange
		}
	}

	if v.SalaryFrom != nil && v.SalaryTo != nil {
		if *v.SalaryFrom > *v.SalaryTo {
			return domain_errors.ErrInvalidSalaryRange
		}
	}

	if !v.IsPaid {
		if v.SalaryFrom != nil || v.SalaryTo != nil {
			return domain_errors.ErrSalaryProvidedForUnpaid
		}
	}

	if v.SalaryTo != nil && *v.SalaryTo > maxSalary {
		return domain_errors.ErrSalaryTooLarge
	}

	if v.SalaryFrom != nil && *v.SalaryFrom < 0 {
		return domain_errors.ErrNegativeSalary
	}

	if v.DurationFromMonths != nil && *v.DurationFromMonths <= 0 {
		return domain_errors.ErrInvalidDurationRange
	}

	if v.HoursPerWeekFrom != nil && *v.HoursPerWeekFrom <= 0 {
		return domain_errors.ErrInvalidHoursRange
	}

	if len(v.Title) == 0 || len(v.Title) > maxTitleLen {
		return domain_errors.ErrInvalidTitleLength
	}

	if len(v.Description) > maxDescriptionLen {
		return domain_errors.ErrInvalidDescriptionLength
	}

	return nil
}
