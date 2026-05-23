package views

import (
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

type VacancySearch struct {
	ID uuid.UUID

	CompanyID   uuid.UUID
	CompanyName string

	Title       string
	Description string

	WorkFormat vacancy.WorkFormat
	City       *string

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

	Status      vacancy.Status
	PublishedAt *time.Time

	ModerationStatus vacancy.ModerationStatus

	CreatedBy uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func SearchViewFromVacancy(vac vacancy.Vacancy, compName string) *VacancySearch {
	return &VacancySearch{
		ID:                vac.ID,
		CompanyID:         vac.CompanyID,
		CompanyName:       compName,
		Title:             vac.Title,
		Description:       vac.Description,
		WorkFormat:        vac.WorkFormat,
		City:              vac.City,
		DurationFromDays:  vac.DurationFromDays,
		DurationToDays:    vac.DurationToDays,
		EmploymentType:    vac.EmploymentType,
		HoursPerWeekFrom:  vac.HoursPerWeekFrom,
		HoursPerWeekTo:    vac.HoursPerWeekTo,
		FlexibleSchedule:  vac.FlexibleSchedule,
		IsPaid:            vac.IsPaid,
		SalaryFrom:        vac.SalaryFrom,
		SalaryTo:          vac.SalaryTo,
		InternshipToOffer: vac.InternshipToOffer,
		Status:            vac.Status,
		PublishedAt:       vac.PublishedAt,
		ModerationStatus:  vac.ModerationStatus,
		CreatedBy:         vac.CreatedBy,
		CreatedAt:         vac.CreatedAt,
		UpdatedAt:         vac.UpdatedAt,
	}
}
