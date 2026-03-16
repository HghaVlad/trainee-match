package get_published_vacancy

import (
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
)

type Response struct {
	ID        uuid.UUID `db:"id"`
	CompanyID uuid.UUID `db:"company_id"`

	CompanyName string `db:"company_name"`

	Title       string                 `db:"title"`
	Description string                 `db:"description"`
	WorkFormat  value_types.WorkFormat `db:"work_format"`
	City        *string                `db:"city"`

	DurationFromDays *int `db:"duration_from_days"`
	DurationToDays   *int `db:"duration_to_days"`

	EmploymentType   value_types.EmploymentType `db:"employment_type"`
	HoursPerWeekFrom *int                       `db:"hours_per_week_from"`
	HoursPerWeekTo   *int                       `db:"hours_per_week_to"`

	FlexibleSchedule bool `db:"flexible_schedule"`

	IsPaid     bool `db:"is_paid"`
	SalaryFrom *int `db:"salary_from"`
	SalaryTo   *int `db:"salary_to"`

	InternshipToOffer bool      `db:"internship_to_offer"`
	PublishedAt       time.Time `db:"published_at"`
}
