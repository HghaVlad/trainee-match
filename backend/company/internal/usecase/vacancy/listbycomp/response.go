package listbycomp

import (
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

type VacancySummary struct {
	ID uuid.UUID `db:"id"`

	Title      string
	WorkFormat vacancy.WorkFormat
	City       *string

	EmploymentType vacancy.EmploymentType

	IsPaid     bool
	SalaryFrom *int
	SalaryTo   *int

	Status    vacancy.Status
	CreatedAt time.Time
}

type Response struct {
	Vacancies  []VacancySummary
	NextCursor *string
}
