package listbycomp

import (
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

type VacancySummary struct {
	ID uuid.UUID `db:"id"`

	Title      string             `db:"title"`
	WorkFormat vacancy.WorkFormat `db:"work_format"`
	City       *string            `db:"city"`

	EmploymentType vacancy.EmploymentType `db:"employment_type"`

	IsPaid     bool `db:"is_paid"`
	SalaryFrom *int `db:"salary_from"`
	SalaryTo   *int `db:"salary_to"`

	PublishedAt time.Time `db:"published_at"`
}

type Response struct {
	Vacancies  []VacancySummary
	NextCursor *string
}
