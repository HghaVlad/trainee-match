package mapper

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/create"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list_by_company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/update"
)

func VacancyToDtoResponse(v *domain.Vacancy) *dto.VacancyFullResponse {
	return &dto.VacancyFullResponse{
		ID:        v.ID,
		CompanyID: v.CompanyID,

		Title:       v.Title,
		Description: v.Description,

		WorkFormat: string(v.WorkFormat),
		City:       v.City,

		DurationFromDays: v.DurationFromDays,
		DurationToDays:   v.DurationToDays,

		EmploymentType:   string(v.EmploymentType),
		HoursPerWeekFrom: v.HoursPerWeekFrom,
		HoursPerWeekTo:   v.HoursPerWeekTo,

		FlexibleSchedule: v.FlexibleSchedule,

		IsPaid:     v.IsPaid,
		SalaryFrom: v.SalaryFrom,
		SalaryTo:   v.SalaryTo,

		InternshipToOffer: v.InternshipToOffer,

		Status:      string(v.Status),
		CreatedBy:   v.CreatedBy,
		PublishedAt: v.PublishedAt,

		CreatedAt:   v.CreatedAt,
		UpdatedAtAt: v.UpdatedAtAt,
	}
}

func VacancyCreateReqToUC(dtoReq *dto.VacancyCreateRequest, companyID uuid.UUID) *create_vacancy.Request {
	req := &create_vacancy.Request{
		CompanyID: companyID,

		Title:       dtoReq.Title,
		Description: dtoReq.Description,

		WorkFormat: value_types.WorkFormat(dtoReq.WorkFormat),
		City:       dtoReq.City,

		DurationFromDays: dtoReq.DurationFromDays,
		DurationToDays:   dtoReq.DurationToDays,

		HoursPerWeekFrom: dtoReq.HoursPerWeekFrom,
		HoursPerWeekTo:   dtoReq.HoursPerWeekTo,

		FlexibleSchedule: dtoReq.FlexibleSchedule,

		IsPaid:     dtoReq.IsPaid,
		SalaryFrom: dtoReq.SalaryFrom,
		SalaryTo:   dtoReq.SalaryTo,

		InternshipToOffer: dtoReq.InternshipToOffer,
	}

	if dtoReq.EmploymentType != nil {
		et := value_types.EmploymentType(*dtoReq.EmploymentType)
		req.EmploymentType = &et
	}

	return req
}

func VacancyListRespToDto(
	resp *list_vacancy.Response,
) *dto.VacancyListResponse {

	items := make([]dto.VacancyListItemResponse, 0, len(resp.Vacancies))

	for _, v := range resp.Vacancies {
		items = append(items, dto.VacancyListItemResponse{
			ID:        v.ID,
			CompanyID: v.CompanyID,

			CompanyName: v.CompanyName,

			Title:      v.Title,
			WorkFormat: string(v.WorkFormat),
			City:       v.City,

			EmploymentType: string(v.EmploymentType),

			IsPaid:     v.IsPaid,
			SalaryFrom: v.SalaryFrom,
			SalaryTo:   v.SalaryTo,

			PublishedAt: v.PublishedAt,
		})
	}

	return &dto.VacancyListResponse{
		Vacancies:  items,
		NextCursor: resp.NextCursor,
	}
}

func ListVacByCompRespToDto(resp *list_vac_by_comp.Response) *dto.VacancyByCompListResponse {
	items := make([]dto.VacancyByCompListItemResponse, 0, len(resp.Vacancies))

	for _, v := range resp.Vacancies {
		items = append(items, dto.VacancyByCompListItemResponse{
			ID: v.ID,

			Title:      v.Title,
			WorkFormat: string(v.WorkFormat),
			City:       v.City,

			EmploymentType: string(v.EmploymentType),

			IsPaid:     v.IsPaid,
			SalaryFrom: v.SalaryFrom,
			SalaryTo:   v.SalaryTo,

			PublishedAt: v.PublishedAt,
		})
	}

	return &dto.VacancyByCompListResponse{
		Vacancies:  items,
		NextCursor: resp.NextCursor,
	}
}

func VacancyCreateRespToDto(resp *create_vacancy.Response) *dto.VacancyCreatedResponse {
	return &dto.VacancyCreatedResponse{
		ID: resp.ID,
	}
}

func VacancyUpdateReqToUC(
	dtoReq *dto.VacancyUpdateRequest,
	companyID uuid.UUID,
	vacancyID uuid.UUID,
) *update_vacancy.Request {

	req := &update_vacancy.Request{
		CompanyID: companyID,
		VacancyID: vacancyID,

		Title:       dtoReq.Title,
		Description: dtoReq.Description,

		City: dtoReq.City,

		DurationFromDays: dtoReq.DurationFromDays,
		DurationToDays:   dtoReq.DurationToDays,

		HoursPerWeekFrom: dtoReq.HoursPerWeekFrom,
		HoursPerWeekTo:   dtoReq.HoursPerWeekTo,

		FlexibleSchedule: dtoReq.FlexibleSchedule,

		IsPaid:     dtoReq.IsPaid,
		SalaryFrom: dtoReq.SalaryFrom,
		SalaryTo:   dtoReq.SalaryTo,

		InternshipToOffer: dtoReq.InternshipToOffer,
	}

	if dtoReq.WorkFormat != nil {
		wf := value_types.WorkFormat(*dtoReq.WorkFormat)
		req.WorkFormat = &wf
	}

	if dtoReq.WorkFormat != nil {
		et := value_types.EmploymentType(*dtoReq.EmploymentType)
		req.EmploymentType = &et
	}

	return req
}
