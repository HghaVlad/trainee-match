package vacancy_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

func validVacancy() *vacancy.Vacancy {
	now := time.Now()

	return &vacancy.Vacancy{
		ID:        uuid.New(),
		CompanyID: uuid.New(),

		Title:       "Backend Intern",
		Description: "Some description",

		WorkFormat:     vacancy.WorkFormatRemote,
		EmploymentType: vacancy.EmploymentTypeFullTime,

		IsPaid: false,

		Status:      vacancy.VacancyStatusPublished,
		PublishedAt: &now,
		CreatedAt:   now,
		UpdatedAt:   time.Now(),
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
		mod  func(v *vacancy.Vacancy)
		err  error
	}{
		{
			name: "invalid work format",
			mod: func(v *vacancy.Vacancy) {
				v.WorkFormat = "invalid"
			},
			err: vacancy.ErrInvalidWorkFormat,
		},
		{
			name: "invalid employment type",
			mod: func(v *vacancy.Vacancy) {
				v.EmploymentType = "invalid"
			},
			err: vacancy.ErrInvalidEmploymentType,
		},
		{
			name: "invalid duration range",
			mod: func(v *vacancy.Vacancy) {
				from, to := 12, 6
				v.DurationFromDays = &from
				v.DurationToDays = &to
			},
			err: vacancy.ErrInvalidDurationRange,
		},
		{
			name: "salary for unpaid vacancy",
			mod: func(v *vacancy.Vacancy) {
				s := 1000
				v.IsPaid = false
				v.SalaryFrom = &s
			},
			err: vacancy.ErrSalaryProvidedForUnpaid,
		},
		{
			name: "salary too large",
			mod: func(v *vacancy.Vacancy) {
				s := 1_000_000_000
				v.IsPaid = true
				v.SalaryTo = &s
			},
			err: vacancy.ErrSalaryTooLarge,
		},
		{
			name: "empty title",
			mod: func(v *vacancy.Vacancy) {
				v.Title = ""
			},
			err: vacancy.ErrInvalidTitleLength,
		},
		{
			name: "description too long",
			mod: func(v *vacancy.Vacancy) {
				v.Description = string(make([]byte, 10000))
			},
			err: vacancy.ErrInvalidDescriptionLength,
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
