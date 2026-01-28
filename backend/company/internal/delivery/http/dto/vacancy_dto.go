package dto

import (
	"time"

	"github.com/google/uuid"
)

type VacancyResponse struct {
	ID        uuid.UUID `json:"id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	CompanyID uuid.UUID `json:"company_id" example:"d290f1ee-6c54-4b01-90e6-d701748f0851"`

	Title       string `json:"title" example:"Go Backend Developer Intern"`
	Description string `json:"description" example:"Join Google's backend team to build scalable services in Go."`

	WorkFormat string  `json:"work_format" enums:"onsite,remote,hybrid" example:"hybrid"`
	City       *string `json:"city,omitempty" example:"Mountain View"`

	DurationFromMonths *int `json:"duration_from_months,omitempty" example:"3"`
	DurationToMonths   *int `json:"duration_to_months,omitempty" example:"6"`

	EmploymentType   string `json:"employment_type" enums:"internship,full_time,part_time" example:"internship"`
	HoursPerWeekFrom *int   `json:"hours_per_week_from,omitempty" example:"30"`
	HoursPerWeekTo   *int   `json:"hours_per_week_to,omitempty" example:"40"`

	FlexibleSchedule bool `json:"flexible_schedule" example:"true"`

	IsPaid     bool `json:"is_paid" example:"true"`
	SalaryFrom *int `json:"salary_from,omitempty" example:"3500"`
	SalaryTo   *int `json:"salary_to,omitempty" example:"5000"`

	InternshipToOffer bool `json:"internship_to_offer" example:"true"`

	IsActive    bool      `json:"is_active" example:"true"`
	PublishedAt time.Time `json:"published_at" example:"2026-01-20T10:00:00Z"`
	CreatedAt   time.Time `json:"created_at" example:"2026-01-18T09:30:00Z"`
	UpdatedAtAt time.Time `json:"updated_at" example:"2026-01-22T14:15:00Z"`
}

type VacancyCreateRequest struct {
	Title       string `json:"title" example:"Go Backend Developer Intern"`
	Description string `json:"description" example:"Join Google's backend team to build scalable services in Go."`

	WorkFormat string  `json:"work_format" enums:"onsite,remote,hybrid" example:"hybrid"`
	City       *string `json:"city,omitempty" example:"Mountain View"`

	DurationFromMonths *int `json:"duration_from_months,omitempty" example:"3"`
	DurationToMonths   *int `json:"duration_to_months,omitempty" example:"6"`

	EmploymentType   *string `json:"employment_type" example:"internship"`
	HoursPerWeekFrom *int    `json:"hours_per_week_from,omitempty" example:"20"`
	HoursPerWeekTo   *int    `json:"hours_per_week_to,omitempty" example:"40"`

	FlexibleSchedule bool `json:"flexible_schedule"`

	IsPaid     bool `json:"is_paid"`
	SalaryFrom *int `json:"salary_from,omitempty" example:"1000"`
	SalaryTo   *int `json:"salary_to,omitempty" example:"1500"`

	InternshipToOffer bool `json:"internship_to_offer"`
}

type VacancyCreatedResponse struct {
	ID uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type VacancyUpdateRequest struct {
	Title       *string `json:"title,omitempty" example:"Go Backend Developer Intern"`
	Description *string `json:"description,omitempty" example:"Work on high-load backend services using Go and PostgreSQL."`

	WorkFormat *string `json:"work_format,omitempty" enums:"onsite,remote,hybrid" example:"remote"`
	City       *string `json:"city,omitempty" example:"Berlin"`

	DurationFromMonths *int `json:"duration_from_months,omitempty" example:"3"`
	DurationToMonths   *int `json:"duration_to_months,omitempty" example:"6"`

	EmploymentType   *string `json:"employment_type,omitempty" enums:"internship,full_time,part_time" example:"internship"`
	HoursPerWeekFrom *int    `json:"hours_per_week_from,omitempty" example:"20"`
	HoursPerWeekTo   *int    `json:"hours_per_week_to,omitempty" example:"40"`

	FlexibleSchedule *bool `json:"flexible_schedule,omitempty" example:"true"`

	IsPaid     *bool `json:"is_paid,omitempty" example:"true"`
	SalaryFrom *int  `json:"salary_from,omitempty" example:"1000"`
	SalaryTo   *int  `json:"salary_to,omitempty" example:"1500"`

	InternshipToOffer *bool `json:"internship_to_offer,omitempty" example:"true"`
}
