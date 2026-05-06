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

	Status      Status
	PublishedAt *time.Time

	CreatedBy uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

const (
	MaxTitleLen       = 200
	MaxDescriptionLen = 5000
)

const (
	MaxSalary       = 15_000_000
	MaxDurationDays = 1800
	MaxHoursPerWeek = 80
)

// Validate checks domain invariants
func (v *Vacancy) Validate() error {
	validators := []func() error{
		v.validateSalary,
		v.validateDuration,
		v.validateHoursPerWeek,
		v.validateWorkFormat,
		v.validateEmploymentType,
		v.validateTitle,
		v.validateDescription,
	}

	for _, validate := range validators {
		if err := validate(); err != nil {
			return err
		}
	}

	if !v.Status.IsValid() {
		return ErrInvalidStatus
	}

	return nil
}

func (v *Vacancy) validateSalary() error {
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

	return nil
}

func (v *Vacancy) validateDuration() error {
	if v.DurationFromDays != nil && v.DurationToDays != nil {
		if *v.DurationFromDays > *v.DurationToDays {
			return ErrInvalidDurationRange
		}
	}

	if v.DurationFromDays != nil && *v.DurationFromDays <= 0 ||
		v.DurationToDays != nil && *v.DurationToDays > MaxDurationDays {
		return ErrInvalidDurationRange
	}

	return nil
}

func (v *Vacancy) validateHoursPerWeek() error {
	if v.HoursPerWeekFrom != nil && v.HoursPerWeekTo != nil {
		if *v.HoursPerWeekFrom > *v.HoursPerWeekTo {
			return ErrInvalidHoursRange
		}
	}

	if v.HoursPerWeekFrom != nil && *v.HoursPerWeekFrom <= 0 ||
		v.HoursPerWeekTo != nil && *v.HoursPerWeekTo > MaxHoursPerWeek {
		return ErrInvalidHoursRange
	}

	return nil
}

func (v *Vacancy) validateWorkFormat() error {
	if !v.WorkFormat.IsValid() {
		return ErrInvalidWorkFormat
	}

	return nil
}

func (v *Vacancy) validateEmploymentType() error {
	if !v.EmploymentType.IsValid() {
		return ErrInvalidEmploymentType
	}

	return nil
}

func (v *Vacancy) validateTitle() error {
	if len([]rune(v.Title)) == 0 || len([]rune(v.Title)) > MaxTitleLen {
		return ErrInvalidTitleLength
	}

	return nil
}

func (v *Vacancy) validateDescription() error {
	if len([]rune(v.Description)) > MaxDescriptionLen {
		return ErrInvalidDescriptionLength
	}

	return nil
}
