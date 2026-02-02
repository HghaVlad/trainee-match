package delete_company

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Usecase struct {
	repo  CompanyRepo
	cache CacheRepo
}

func NewUsecase(
	repo CompanyRepo,
	cache CacheRepo,
) *Usecase {
	return &Usecase{
		repo:  repo,
		cache: cache,
	}
}

func (u *Usecase) Execute(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := u.repo.Delete(ctx, id)

	if err != nil {
		return err
	}

	u.cache.Del(ctx, id)
	return nil
}
