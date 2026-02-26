package create_company

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
)

type Usecase struct {
	compRepo   CompanyRepo
	memberRepo CompanyMemberRepo
	txManager  uc_common.TxManager
}

func NewUsecase(compRepo CompanyRepo, memberRepo CompanyMemberRepo, txManager uc_common.TxManager) *Usecase {
	return &Usecase{
		compRepo:   compRepo,
		memberRepo: memberRepo,
		txManager:  txManager,
	}
}

func (u *Usecase) Execute(ctx context.Context, request *Request, identity uc_common.Identity) (*Response, error) {
	if identity.Role != uc_common.RoleHR {
		return nil, domain_errors.ErrHrRoleRequired
	}

	company := &domain.Company{
		ID:          uuid.New(),
		Name:        request.Name,
		Description: request.Description,
		Website:     request.Website,
	}

	valErr := company.Validate()
	if valErr != nil {
		return nil, valErr
	}

	// creator is an admin of company
	member := &domain.CompanyMember{
		UserID:    identity.UserID,
		CompanyID: company.ID,
		Role:      value_types.CompanyRoleAdmin,
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		err := u.compRepo.Create(ctx, company)
		if err != nil {
			return err
		}

		return u.memberRepo.Create(ctx, member)
	})

	if err != nil {
		return nil, err
	}

	resp := &Response{
		ID: company.ID,
	}

	return resp, nil
}
