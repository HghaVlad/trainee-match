package add_member

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
)

type Usecase struct {
	memberRepo CompanyMemberRepo
}

func NewUsecase(memberRepo CompanyMemberRepo) *Usecase {
	return &Usecase{memberRepo: memberRepo}
}

func (u *Usecase) Execute(ctx context.Context, req *Request, identity uc_common.Identity) error {
	if err := req.Validate(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := u.authorize(ctx, req.CompanyID, identity); err != nil {
		return err
	}

	member := &domain.CompanyMember{
		UserID:    req.UserID,
		CompanyID: req.CompanyID,
		Role:      req.Role,
	}

	return u.memberRepo.Create(ctx, member)
}

func (u *Usecase) authorize(ctx context.Context, companyID uuid.UUID, identity uc_common.Identity) error {
	if identity.Role != uc_common.RoleHR {
		return domain_errors.ErrHrRoleRequired
	}

	member, err := u.memberRepo.Get(ctx, identity.UserID, companyID)
	if errors.Is(err, domain_errors.ErrCompanyMemberNotFound) {
		return domain_errors.ErrCompanyMemberRequired
	}
	if err != nil {
		return err
	}

	if member.Role != value_types.CompanyRoleAdmin {
		return domain_errors.ErrInsufficientRoleInCompany
	}

	return nil
}
