package delete_member

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

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

func (u *Usecase) Execute(ctx context.Context, companyID, userID uuid.UUID, identity uc_common.Identity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := u.authorize(ctx, companyID, identity); err != nil {
		return err
	}

	return u.memberRepo.Delete(ctx, userID, companyID)
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
