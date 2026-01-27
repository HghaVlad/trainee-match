package domain

import (
	"time"

	"github.com/google/uuid"
)

type Vacancy struct {
	ID        uuid.UUID `db:"id"`
	CompanyID uuid.UUID `db:"company_id"`

	Title       string `db:"title"`
	Description string `db:"description"`

	WorkFormat WorkFormat `db:"work_format"`
	City       *string    `db:"city"`

	DurationFromMonths *int `db:"duration_from_months"`
	DurationToMonths   *int `db:"duration_to_months"`

	EmploymentType   EmploymentType `db:"employment_type"`
	HoursPerWeekFrom *int           `db:"hours_per_week_from"`
	HoursPerWeekTo   *int           `db:"hours_per_week_to"`

	FlexibleSchedule bool `db:"flexible_schedule"`

	IsPaid     bool `db:"is_paid"`
	SalaryFrom *int `db:"salary_from"`
	SalaryTo   *int `db:"salary_to"`

	InternshipToOffer bool `db:"internship_to_offer"`

	IsActive    bool      `db:"is_active"`
	PublishedAt time.Time `db:"published_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAtAt time.Time `db:"updated_at"`
}

type WorkFormat string

const (
	WorkFormatOnSite WorkFormat = "onsite"
	WorkFormatRemote WorkFormat = "remote"
	WorkFormatHybrid WorkFormat = "hybrid"
)

type EmploymentType string

const (
	EmploymentTypeInternship EmploymentType = "internship'"
	EmploymentTypeFullTime   EmploymentType = "full_time"
	EmploymentTypePartTime   EmploymentType = "part_time"
)
