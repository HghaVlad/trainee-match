package dto

import (
	"time"

	"github.com/google/uuid"
)

type VacancyFullResponse struct {
	ID        uuid.UUID `json:"id"        example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	CompanyID uuid.UUID `json:"companyId" example:"d290f1ee-6c54-4b01-90e6-d701748f0851"`

	Title       string `json:"title"       example:"Go Backend Developer Intern"`
	Description string `json:"description" example:"Join Google's backend team to build scalable services in Go."`

	WorkFormat string  `json:"workFormat"     enums:"onsite,remote,hybrid" example:"hybrid"`
	City       *string `json:"city,omitempty"                              example:"Mountain View"`

	DurationFromDays *int `json:"durationFromDays,omitempty" example:"60"`
	DurationToDays   *int `json:"durationToDays,omitempty"   example:"90"`

	EmploymentType   string `json:"employmentType"             enums:"internship,full_time,part_time" example:"internship"`
	HoursPerWeekFrom *int   `json:"hoursPerWeekFrom,omitempty"                                        example:"30"`
	HoursPerWeekTo   *int   `json:"hoursPerWeekTo,omitempty"                                          example:"40"`

	FlexibleSchedule bool `json:"flexibleSchedule" example:"true"`

	IsPaid     bool `json:"isPaid"               example:"true"`
	SalaryFrom *int `json:"salaryFrom,omitempty" example:"3500"`
	SalaryTo   *int `json:"salaryTo,omitempty"   example:"5000"`

	InternshipToOffer bool `json:"internshipToOffer" example:"true"`

	Status    string    `json:"status"    enums:"draft,published,archived" example:"published"`
	CreatedBy uuid.UUID `json:"createdBy"                                  example:"d290f1ee-6c54-4b01-90e6-d701748f0851"`

	PublishedAt *time.Time `json:"publishedAt,omitempty" example:"2026-01-20T10:00:00Z"`

	CreatedAt   time.Time `json:"createdAt" example:"2026-01-18T09:30:00Z"`
	UpdatedAtAt time.Time `json:"updatedAt" example:"2026-01-22T14:15:00Z"`
}

type VacancyPublicResponse struct {
	ID        uuid.UUID `json:"id"        example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	CompanyID uuid.UUID `json:"companyId" example:"d290f1ee-6c54-4b01-90e6-d701748f0851"`

	CompanyName string `json:"companyName" example:"Google Inc."`

	Title       string `json:"title"       example:"Go Backend Developer Intern"`
	Description string `json:"description" example:"Join Google's backend team to build scalable services in Go."`

	WorkFormat string  `json:"workFormat"     enums:"onsite,remote,hybrid" example:"hybrid"`
	City       *string `json:"city,omitempty"                              example:"Mountain View"`

	DurationFromDays *int `json:"durationFromDays,omitempty" example:"60"`
	DurationToDays   *int `json:"durationToDays,omitempty"   example:"90"`

	EmploymentType   string `json:"employmentType"             enums:"internship,full_time,part_time" example:"internship"`
	HoursPerWeekFrom *int   `json:"hoursPerWeekFrom,omitempty"                                        example:"30"`
	HoursPerWeekTo   *int   `json:"hoursPerWeekTo,omitempty"                                          example:"40"`

	FlexibleSchedule bool `json:"flexibleSchedule" example:"true"`

	IsPaid     bool `json:"isPaid"               example:"true"`
	SalaryFrom *int `json:"salaryFrom,omitempty" example:"3500"`
	SalaryTo   *int `json:"salaryTo,omitempty"   example:"5000"`

	InternshipToOffer bool      `json:"internshipToOffer" example:"true"`
	PublishedAt       time.Time `json:"publishedAt"       example:"2026-01-20T10:00:00Z"`
}

type VacancyListItemResponse struct {
	ID        uuid.UUID `json:"id"        example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	CompanyID uuid.UUID `json:"companyId" example:"d290f1ee-6c54-4b01-90e6-d701748f0851"`

	CompanyName string `json:"companyName" example:"Google Inc."`

	Title      string  `json:"title"          example:"Go Backend Developer Intern"`
	WorkFormat string  `json:"workFormat"     example:"hybrid"`
	City       *string `json:"city,omitempty" example:"Mountain View"`

	EmploymentType string `json:"employmentType" example:"internship,full_time,part_time"`

	IsPaid     bool `json:"isPaid"               example:"true"`
	SalaryFrom *int `json:"salaryFrom,omitempty" example:"3500"`
	SalaryTo   *int `json:"salaryTo,omitempty"   example:"5000"`

	PublishedAt time.Time `json:"publishedAt" example:"2026-01-18T09:30:00Z"`
}

type VacancyListResponse struct {
	Vacancies  []VacancyListItemResponse `json:"vacancies"`
	NextCursor *string                   `json:"nextCursor,omitempty"`
}

type VacancyByCompListItemResponse struct {
	ID uuid.UUID `json:"id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`

	Title      string  `json:"title"          example:"Go Backend Developer Intern"`
	WorkFormat string  `json:"workFormat"     example:"hybrid"`
	City       *string `json:"city,omitempty" example:"Mountain View"`

	EmploymentType string `json:"employmentType" example:"internship,full_time,part_time"`

	IsPaid     bool `json:"isPaid"               example:"true"`
	SalaryFrom *int `json:"salaryFrom,omitempty" example:"3500"`
	SalaryTo   *int `json:"salaryTo,omitempty"   example:"5000"`

	PublishedAt time.Time `json:"publishedAt" example:"2026-01-18T09:30:00Z"`
}

type VacancyByCompListResponse struct {
	Vacancies  []VacancyByCompListItemResponse `json:"vacancies"`
	NextCursor *string                         `json:"nextCursor,omitempty"`
}

type VacancyCreateRequest struct {
	Title       string `json:"title"       example:"Go Backend Developer Intern"`
	Description string `json:"description" example:"Join Google's backend team to build scalable services in Go."`

	WorkFormat string  `json:"workFormat"     enums:"onsite,remote,hybrid" example:"hybrid"`
	City       *string `json:"city,omitempty"                              example:"Mountain View"`

	DurationFromDays *int `json:"durationFromDays,omitempty" example:"60"`
	DurationToDays   *int `json:"durationToDays,omitempty"   example:"90"`

	EmploymentType   *string `json:"employmentType"             example:"internship"`
	HoursPerWeekFrom *int    `json:"hoursPerWeekFrom,omitempty" example:"20"`
	HoursPerWeekTo   *int    `json:"hoursPerWeekTo,omitempty"   example:"40"`

	FlexibleSchedule bool `json:"flexibleSchedule"`

	IsPaid     bool `json:"isPaid"`
	SalaryFrom *int `json:"salaryFrom,omitempty" example:"1000"`
	SalaryTo   *int `json:"salaryTo,omitempty"   example:"1500"`

	InternshipToOffer bool `json:"internshipToOffer"`
}

type VacancyCreatedResponse struct {
	ID uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type VacancyUpdateRequest struct {
	Title       *string `json:"title,omitempty"       example:"Go Backend Developer Intern"`
	Description *string `json:"description,omitempty" example:"Work on high-load backend services using Go and PostgreSQL."`

	WorkFormat *string `json:"workFormat,omitempty" enums:"onsite,remote,hybrid" example:"remote"`
	City       *string `json:"city,omitempty"                                    example:"Berlin"`

	DurationFromDays *int `json:"durationFromDays,omitempty" example:"3"`
	DurationToDays   *int `json:"durationToDays,omitempty"   example:"6"`

	EmploymentType   *string `json:"employmentType,omitempty"   enums:"internship,full_time,part_time" example:"internship"`
	HoursPerWeekFrom *int    `json:"hoursPerWeekFrom,omitempty"                                        example:"20"`
	HoursPerWeekTo   *int    `json:"hoursPerWeekTo,omitempty"                                          example:"40"`

	FlexibleSchedule *bool `json:"flexibleSchedule,omitempty" example:"true"`

	IsPaid     *bool `json:"isPaid,omitempty"     example:"true"`
	SalaryFrom *int  `json:"salaryFrom,omitempty" example:"1000"`
	SalaryTo   *int  `json:"salaryTo,omitempty"   example:"1500"`

	InternshipToOffer *bool `json:"internshipToOffer,omitempty" example:"true"`
}
