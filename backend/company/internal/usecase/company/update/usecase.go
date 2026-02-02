package update_company

import (
	"context"
	"time"
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

func (u *Usecase) Execute(ctx context.Context, req *Request) error {
	if err := req.Validate(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := u.repo.Update(ctx, req)

	if err != nil {
		return err
	}

	u.cache.Del(ctx, req.ID)
	return nil
}
