package create_company

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	uc_common "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
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

func (u *Usecase) Execute(ctx context.Context, request *Request) (*Response, error) {

	// TODO: do smth with owner id

	company := &entities.Company{
		ID:          uuid.New(),
		Name:        request.Name,
		Description: request.Description,
		Website:     request.Website,
		OwnerID:     uuid.New(),
	}

	err := u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		return u.repo.Create(ctx, company)
	})

	if err != nil {
		return nil, err
	}

	resp := &Response{
		ID: company.ID,
	}

	return resp, nil
}
