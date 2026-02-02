package delete_vacancy

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Usecase struct {
	repo  VacancyRepo
	cache CacheRepo
}

func NewUsecase(repo VacancyRepo, cache CacheRepo) *Usecase {
	return &Usecase{repo: repo, cache: cache}
}

func (u *Usecase) Execute(ctx context.Context, vacancyID uuid.UUID, companyID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := u.repo.Delete(ctx, vacancyID, companyID)
	if err != nil {
		return err
	}

	u.cache.Del(ctx, vacancyID)
	return nil
}
