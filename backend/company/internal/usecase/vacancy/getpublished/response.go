package getpublished

import (
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

type Response struct {
	ID        uuid.UUID
	CompanyID uuid.UUID

	CompanyName string

	Title       string
	Description string
	WorkFormat  vacancy.WorkFormat
	City        *string

	DurationFromDays *int
	DurationToDays   *int

	EmploymentType   vacancy.EmploymentType
	HoursPerWeekFrom *int
	HoursPerWeekTo   *int

	FlexibleSchedule bool

	IsPaid     bool
	SalaryFrom *int
	SalaryTo   *int

	InternshipToOffer bool
	PublishedAt       time.Time
}
