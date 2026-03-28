package create

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
)

// Usecase of company creation.
// Adds creator as an admin member of company.
type Usecase struct {
	compRepo   CompanyRepo
	memberRepo CompanyMemberRepo
	txManager  common.TxManager
}

func NewUsecase(compRepo CompanyRepo, memberRepo CompanyMemberRepo, txManager common.TxManager) *Usecase {
	return &Usecase{
		compRepo:   compRepo,
		memberRepo: memberRepo,
		txManager:  txManager,
	}
}

// Execute create company.
// Adds creator as an admin member of company.
func (u *Usecase) Execute(ctx context.Context, request *Request, ident *identity.Identity) (*Response, error) {
	if ident.Role != identity.RoleHR {
		return nil, identity.ErrHrRoleRequired
	}

	comp := &company.Company{
		ID:          uuid.New(),
		Name:        request.Name,
		Description: request.Description,
		Website:     request.Website,
	}

	valErr := comp.Validate()
	if valErr != nil {
		return nil, valErr
	}

	// creator is an admin of company
	memb := &member.CompanyMember{
		UserID:    ident.UserID,
		CompanyID: comp.ID,
		Role:      member.CompanyRoleAdmin,
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		err := u.compRepo.Create(ctx, comp)
		if err != nil {
			return err
		}

		return u.memberRepo.Create(ctx, memb)
	})

	if err != nil {
		return nil, err
	}

	resp := &Response{
		ID: comp.ID,
	}

	return resp, nil
}
