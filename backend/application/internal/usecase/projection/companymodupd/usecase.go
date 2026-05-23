package companymodupd

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type Usecase struct {
	repo vacancyRepo
}

func NewUsecase(repo vacancyRepo) Usecase {
	return Usecase{
		repo: repo,
	}
}

func (u Usecase) Execute(ctx context.Context, event projection.CompanyModerationUpdatedEvent) error {
	return u.repo.UpdateCompanyModStatus(ctx, event.CompanyID, event.ModStatus)
}
