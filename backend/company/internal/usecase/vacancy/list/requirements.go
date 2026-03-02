package list_vacancy

import (
	"errors"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
	"github.com/google/uuid"
)

type RangeInt struct {
	Min *int
	Max *int
}

type Requirements struct {
	Salary            *RangeInt
	HoursPerWeek      *RangeInt
	Duration          *RangeInt
	WorkFormat        *[]value_types.WorkFormat
	Companies         *[]uuid.UUID
	City              *[]string
	IsPaid            *bool
	InternshipToOffer *bool
	FlexibleSchedule  *bool
}

func (r *Requirements) Validate() error {
	if r.Salary != nil {
		if err := validateRange(r.Salary, 0, domain.MaxSalary); err != nil {
			return domain_errors.ErrInvalidSalaryRange
		}
	}

	if r.HoursPerWeek != nil {
		if err := validateRange(r.HoursPerWeek, 1, domain.MaxHoursPerWeek); err != nil {
			return domain_errors.ErrInvalidHoursRange
		}
	}

	if r.Duration != nil {
		if err := validateRange(r.Duration, 1, domain.MaxDurationDays); err != nil {
			return domain_errors.ErrInvalidDurationRange
		}
	}

	if r.WorkFormat != nil {
		for _, wf := range *r.WorkFormat {
			if !wf.IsValid() {
				return domain_errors.ErrInvalidWorkFormat
			}
		}
	}

	if r.Companies != nil && len(*r.Companies) == 0 {
		return domain_errors.ErrEmptyCompaniesFilter
	}

	if r.City != nil && len(*r.City) == 0 {
		return domain_errors.ErrEmptyCityFilter
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
