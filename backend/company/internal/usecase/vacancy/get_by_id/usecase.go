package get_vacancy

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
)

type Usecase struct {
	repo  Repository
	cache CacheRepo
}

func NewUsecase(repo Repository, cache CacheRepo) *Usecase {
	return &Usecase{
		repo:  repo,
		cache: cache,
	}
}

func (u *Usecase) Execute(ctx context.Context, vacancyID uuid.UUID, companyID uuid.UUID) (*domain.Vacancy, error) {
	vacancy := u.cache.Get(ctx, vacancyID)

	if vacancy != nil && vacancy.CompanyID == companyID {
		return vacancy, nil
	}

	vacancy, err := u.repo.GetByID(ctx, vacancyID, companyID)
	if err != nil {
		return nil, err
	}

	u.cache.Put(ctx, vacancyID, vacancy, time.Second*300)
	return vacancy, nil
}
