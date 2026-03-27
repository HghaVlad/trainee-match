package list

import (
	"errors"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

type RangeInt struct {
	Min *int
	Max *int
}

type Requirements struct {
	Salary            *RangeInt
	HoursPerWeek      *RangeInt
	Duration          *RangeInt
	WorkFormat        *[]vacancy.WorkFormat
	Companies         *[]uuid.UUID
	City              *[]string
	IsPaid            *bool
	InternshipToOffer *bool
	FlexibleSchedule  *bool
}

func (r *Requirements) Validate() error {
	if r.Salary != nil {
		if err := validateRange(r.Salary, 0, vacancy.MaxSalary); err != nil {
			return vacancy.ErrInvalidSalaryRange
		}
	}

	if r.HoursPerWeek != nil {
		if err := validateRange(r.HoursPerWeek, 1, vacancy.MaxHoursPerWeek); err != nil {
			return vacancy.ErrInvalidHoursRange
		}
	}

	if r.Duration != nil {
		if err := validateRange(r.Duration, 1, vacancy.MaxDurationDays); err != nil {
			return vacancy.ErrInvalidDurationRange
		}
	}

	if r.WorkFormat != nil {
		for _, wf := range *r.WorkFormat {
			if !wf.IsValid() {
				return vacancy.ErrInvalidWorkFormat
			}
		}
	}

	if r.Companies != nil && len(*r.Companies) == 0 {
		return vacancy.ErrEmptyCompaniesFilter
	}

	if r.City != nil && len(*r.City) == 0 {
		return vacancy.ErrEmptyCityFilter
	}

	return nil
}

func validateRange(r *RangeInt, minAllowed, maxAllowed int) error {
	if r.Min != nil && r.Max != nil {
		if *r.Min > *r.Max {
			return errors.New("min > max")
		}
	}
	if r.Min != nil && *r.Min < minAllowed {
		return errors.New("min too small")
	}
	if r.Max != nil && *r.Max > maxAllowed {
		return errors.New("max too large")
	}
	return nil
}
