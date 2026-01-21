package mapper

import (
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/create_company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/get_company"
)

func GetCompRespToDto(company *get_company.Response) *dto.CompanyResponse {
	return &dto.CompanyResponse{
		ID:          company.ID,
		Name:        company.Name,
		Description: company.Description,
		Website:     company.Website,
		OwnerId:     company.OwnerID,
		LogoURL:     company.LogoURL,
		CreatedAt:   company.CreatedAt,
		UpdatedAt:   company.UpdatedAt,
	}
}

func CompanyCreateReqToUC(dtoReq *dto.CompanyCreateRequest) *create_company.Request {
	return &create_company.Request{
		Name:        dtoReq.Name,
		Description: dtoReq.Description,
		Website:     dtoReq.Website,
	}
}

func CompanyCreateRespToDto(resp *create_company.Response) *dto.CompanyCreatedResponse {
	return &dto.CompanyCreatedResponse{
		ID: resp.ID,
	}
}
