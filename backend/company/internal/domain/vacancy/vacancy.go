package vacancy

import (
	"time"

	"github.com/google/uuid"
)

type Vacancy struct {
	ID        uuid.UUID `db:"id"`
	CompanyID uuid.UUID `db:"company_id"`

	Title       string `db:"title"`
	Description string `db:"description"`

	WorkFormat WorkFormat `db:"work_format"`
	City       *string    `db:"city"`

	DurationFromDays *int `db:"duration_from_days"`
	DurationToDays   *int `db:"duration_to_days"`

	EmploymentType   EmploymentType `db:"employment_type"`
	HoursPerWeekFrom *int           `db:"hours_per_week_from"`
	HoursPerWeekTo   *int           `db:"hours_per_week_to"`

	FlexibleSchedule bool `db:"flexible_schedule"`

	IsPaid     bool `db:"is_paid"`
	SalaryFrom *int `db:"salary_from"`
	SalaryTo   *int `db:"salary_to"`

	InternshipToOffer bool `db:"internship_to_offer"`

	Status      VacancyStatus `db:"status"`
	PublishedAt *time.Time    `db:"published_at"`

	CreatedBy   uuid.UUID `db:"created_by_user_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAtAt time.Time `db:"updated_at"`
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
