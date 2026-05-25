package indexsearch

import (
	"context"

	"github.com/google/uuid"
)

type Usecase struct {
	vacRepo    VacancyRepo
	searchRepo SearchRepo
}

func NewUsecase(vacRepo VacancyRepo, searchRepo SearchRepo) *Usecase {
	return &Usecase{
		vacRepo:    vacRepo,
		searchRepo: searchRepo,
	}
}

func (u *Usecase) Index(ctx context.Context, vacID uuid.UUID) error {
	vacSearchView, err := u.vacRepo.GetSearchView(ctx, vacID)
	if err != nil {
		return err
	}

	return u.searchRepo.Index(ctx, *vacSearchView)
}

func (u *Usecase) UpdateCompanyName(ctx context.Context, compID uuid.UUID, newName string) error {
	return u.searchRepo.UpdateCompanyName(ctx, compID, newName)
}

func (u *Usecase) RemoveByCompany(ctx context.Context, companyID uuid.UUID) error {
	return u.searchRepo.RemoveByCompanyID(ctx, companyID)
}
