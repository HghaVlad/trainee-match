package mapper

import (
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/get_company"
)

func GetRespToDto(company *get_company.Response) *dto.CompanyResponse {
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
