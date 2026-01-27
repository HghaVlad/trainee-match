package mapper

import (
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
)

func VacancyToDtoResponse(v *domain.Vacancy) *dto.VacancyResponse {
	return &dto.VacancyResponse{
		ID:        v.ID,
		CompanyID: v.CompanyID,

		Title:       v.Title,
		Description: v.Description,

		WorkFormat: string(v.WorkFormat),
		City:       v.City,

		DurationFromMonths: v.DurationFromMonths,
		DurationToMonths:   v.DurationToMonths,

		EmploymentType:   string(v.EmploymentType),
		HoursPerWeekFrom: v.HoursPerWeekFrom,
		HoursPerWeekTo:   v.HoursPerWeekTo,

		FlexibleSchedule: v.FlexibleSchedule,

		IsPaid:     v.IsPaid,
		SalaryFrom: v.SalaryFrom,
		SalaryTo:   v.SalaryTo,

		InternshipToOffer: v.InternshipToOffer,

		IsActive:    v.IsActive,
		PublishedAt: v.PublishedAt,
		CreatedAt:   v.CreatedAt,
		UpdatedAtAt: v.UpdatedAtAt,
	}
}
