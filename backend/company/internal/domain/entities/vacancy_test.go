package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
)

func validVacancy() *domain.Vacancy {
	now := time.Now()

	return &domain.Vacancy{
		ID:        uuid.New(),
		CompanyID: uuid.New(),

		Title:       "Backend Intern",
		Description: "Some description",

		WorkFormat:     value_types.WorkFormatRemote,
		EmploymentType: value_types.EmploymentTypeFullTime,

		IsPaid: false,

		Status:      value_types.VacancyStatusPublished,
		PublishedAt: &now,
		CreatedAt:   now,
		UpdatedAtAt: time.Now(),
	}
}

func TestVacancy_Validate_OK(t *testing.T) {
	v := validVacancy()
	err := v.Validate()
	require.NoError(t, err)
}

func TestVacancy_Validate_Errors(t *testing.T) {
	tests := []struct {
		name string
		mod  func(v *domain.Vacancy)
		err  error
	}{
		{
			name: "invalid work format",
			mod: func(v *domain.Vacancy) {
				v.WorkFormat = "invalid"
			},
			err: domain_errors.ErrInvalidWorkFormat,
		},
		{
			name: "invalid employment type",
			mod: func(v *domain.Vacancy) {
				v.EmploymentType = "invalid"
			},
			err: domain_errors.ErrInvalidEmploymentType,
		},
		{
			name: "invalid duration range",
			mod: func(v *domain.Vacancy) {
				from, to := 12, 6
				v.DurationFromDays = &from
				v.DurationToDays = &to
			},
			err: domain_errors.ErrInvalidDurationRange,
		},
		{
			name: "salary for unpaid vacancy",
			mod: func(v *domain.Vacancy) {
				s := 1000
				v.IsPaid = false
				v.SalaryFrom = &s
			},
			err: domain_errors.ErrSalaryProvidedForUnpaid,
		},
		{
			name: "salary too large",
			mod: func(v *domain.Vacancy) {
				s := 1_000_000_000
				v.IsPaid = true
				v.SalaryTo = &s
			},
			err: domain_errors.ErrSalaryTooLarge,
		},
		{
			name: "empty title",
			mod: func(v *domain.Vacancy) {
				v.Title = ""
			},
			err: domain_errors.ErrInvalidTitleLength,
		},
		{
			name: "description too long",
			mod: func(v *domain.Vacancy) {
				v.Description = string(make([]byte, 10000))
			},
			err: domain_errors.ErrInvalidDescriptionLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validVacancy()
			tt.mod(v)

			err := v.Validate()

			require.ErrorIs(t, err, tt.err)
		})
	}
}
