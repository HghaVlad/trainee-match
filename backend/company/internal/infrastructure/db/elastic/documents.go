package elastic

import "time"

type VacancyDocument struct {
	ID string `json:"id"`

	CompanyID   string `json:"company_id"`
	CompanyName string `json:"company_name"`

	Title       string `json:"title"`
	Description string `json:"description"`

	WorkFormat string  `json:"work_format"`
	City       *string `json:"city"`

	EmploymentType string `json:"employment_type"`

	DurationFromDays *int `json:"duration_from_days"`
	DurationToDays   *int `json:"duration_to_days"`

	HoursPerWeekFrom *int `json:"hours_per_week_from"`
	HoursPerWeekTo   *int `json:"hours_per_week_to"`

	FlexibleSchedule bool `json:"flexible_schedule"`

	IsPaid bool `json:"is_paid"`

	SalaryFrom *int `json:"salary_from"`
	SalaryTo   *int `json:"salary_to"`

	InternshipToOffer bool `json:"internship_to_offer"`

	Status           string `json:"status"`
	ModerationStatus string `json:"moderation_status"`

	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
}
