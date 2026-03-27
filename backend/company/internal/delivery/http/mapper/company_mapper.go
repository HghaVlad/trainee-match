package mapper

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/create"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/get"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/list"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/update"
)

func GetCompRespToDto(company *get.Response) *dto.CompanyResponse {
	return &dto.CompanyResponse{
		ID:               company.ID,
		Name:             company.Name,
		OpenVacanciesCnt: company.OpenVacanciesCnt,
		Description:      company.Description,
		Website:          company.Website,
		LogoURL:          company.LogoURL,
		CreatedAt:        company.CreatedAt,
		UpdatedAt:        company.UpdatedAt,
	}
}

func CompanyListRespToDto(
	resp *list.Response,
) *dto.CompanyListResponse {

	items := make([]dto.CompanyListItemResponse, 0, len(resp.Companies))

	for _, c := range resp.Companies {
		items = append(items, dto.CompanyListItemResponse{
			ID:               c.ID,
			Name:             c.Name,
			OpenVacanciesCnt: c.OpenVacanciesCnt,
			LogoURL:          c.LogoKey,
		})
	}

	return &dto.CompanyListResponse{
		Companies:  items,
		NextCursor: resp.NextCursor,
	}
}

func CompanyCreateReqToUC(dtoReq *dto.CompanyCreateRequest) *create.Request {
	return &create.Request{
		Name:        dtoReq.Name,
		Description: dtoReq.Description,
		Website:     dtoReq.Website,
	}
}

func CompanyCreateRespToDto(resp *create.Response) *dto.CompanyCreatedResponse {
	return &dto.CompanyCreatedResponse{
		ID: resp.ID,
	}
}

func CompanyUpdateReqToUC(
	id uuid.UUID,
	dtoReq *dto.CompanyUpdateRequest,
) *update.Request {
	return &update.Request{
		ID:          id,
		Name:        dtoReq.Name,
		Description: dtoReq.Description,
		Website:     dtoReq.Website,
	}
}
