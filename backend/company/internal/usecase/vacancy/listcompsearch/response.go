package listcompsearch

import (
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/views"
)

type VacancySummary struct {
	ID uuid.UUID

	Title      string
	WorkFormat vacancy.WorkFormat
	City       *string

	EmploymentType vacancy.EmploymentType

	IsPaid     bool
	SalaryFrom *int
	SalaryTo   *int

	Status    vacancy.Status
	ModStatus vacancy.ModerationStatus
	CreatedAt time.Time
}

type Response struct {
	Vacancies  []views.MemberVacSummary
	NextCursor *string
}
