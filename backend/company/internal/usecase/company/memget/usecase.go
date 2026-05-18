package memget

import (
	"context"
	"time"

	"github.com/google/uuid"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
)

type Usecase struct {
	repo companyRepo
}

func NewUsecase(repo companyRepo) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) Execute(ctx context.Context, id uuid.UUID, iden *identity.Identity) (*Response, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	comp, err := u.repo.GetByMember(ctx, id, iden.UserID)
	if err != nil {
		return nil, err
	}

	resp := toResponse(comp, comp.LogoKey)
	return resp, nil
}

func toResponse(company *domain.Company, logoURL *string) *Response {
	return &Response{
		ID:               company.ID,
		Name:             company.Name,
		OpenVacanciesCnt: company.OpenVacanciesCnt,
		Description:      company.Description,
		Website:          company.Website,
		LogoURL:          logoURL,
		ModStatus:        company.ModerationStatus,
		CreatedAt:        company.CreatedAt,
		UpdatedAt:        company.UpdatedAt,
	}
}
