package get_company

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
)

type GetByIDUsecase struct {
	repo  CompanyRepo
	cache CacheRepo
}

func NewGetByIDUsecase(repo CompanyRepo, cache CacheRepo) *GetByIDUsecase {
	return &GetByIDUsecase{
		repo:  repo,
		cache: cache,
	}
}

func (u *GetByIDUsecase) Execute(ctx context.Context, id uuid.UUID) (*Response, error) {
	// TODO: think about retrieving logo from minio (via presigned or nah)
	company := u.cache.Get(ctx, id)

	if company != nil {
		resp := toResponse(company, company.LogoKey)
		return resp, nil
	}

	company, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	u.cache.Put(ctx, id, company, time.Second*300)

	resp := toResponse(company, company.LogoKey)
	return resp, nil
}

func toResponse(company *domain.Company, logoURL *string) *Response {
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
