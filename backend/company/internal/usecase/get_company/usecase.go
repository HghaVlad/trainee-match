package get_company

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
)

type GetByIDUsecase struct {
	repo CompanyRepo
}

func NewGetByIDUsecase(repo CompanyRepo) *GetByIDUsecase {
	return &GetByIDUsecase{repo: repo}
}

func (u *GetByIDUsecase) Execute(ctx context.Context, id uuid.UUID) (*Response, error) {

	company, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// TODO: think about retrieving logo from minio (via presigned or nah)
	resp := toResponse(company, company.LogoKey)

	return resp, nil
}

func toResponse(company *entities.Company, logoURL *string) *Response {
	return &Response{
		ID:               company.ID,
		Name:             company.Name,
		OpenVacanciesCnt: company.OpenVacanciesCnt,
		Description:      company.Description,
		Website:          company.Website,
		OwnerID:          company.OwnerID,
		LogoURL:          logoURL,
		CreatedAt:        company.CreatedAt,
		UpdatedAt:        company.UpdatedAt,
	}
}
