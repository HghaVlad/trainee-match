package delete_company

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
)

type Usecase struct {
	repo      CompanyRepo
	txManager uc_common.TxManager
}

func NewUsecase(repo CompanyRepo, txManager uc_common.TxManager) *Usecase {
	return &Usecase{
		repo:      repo,
		txManager: txManager,
	}
}

func (u *Usecase) Execute(ctx context.Context, id uuid.UUID) error {
	return u.txManager.WithinTx(ctx, func(ctx context.Context) error {

		return u.repo.Delete(ctx, id)
	})
}
