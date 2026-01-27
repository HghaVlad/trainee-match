package update_company

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
)

type Usecase struct {
	repo      CompanyRepo
	cache     CacheRepo
	txManager uc_common.TxManager
}

func NewUsecase(
	repo CompanyRepo,
	cache CacheRepo,
	txManager uc_common.TxManager,
) *Usecase {
	return &Usecase{
		repo:      repo,
		cache:     cache,
		txManager: txManager,
	}
}

func (u *Usecase) Execute(ctx context.Context, req *Request) error {
	err := u.txManager.WithinTx(ctx, func(ctx context.Context) error {

		return u.repo.Update(ctx, req)
	})

	if err != nil {
		return err
	}

	u.cache.Del(ctx, req.ID)
	return nil
}
