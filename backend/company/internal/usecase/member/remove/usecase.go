package remove

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
)

type Usecase struct {
	memberRepo CompanyMemberRepo
}

func NewUsecase(memberRepo CompanyMemberRepo) *Usecase {
	return &Usecase{memberRepo: memberRepo}
}

func (u *Usecase) Execute(ctx context.Context, companyID, userID uuid.UUID, identity identity.Identity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := u.authorize(ctx, companyID, identity); err != nil {
		return err
	}

	return u.memberRepo.Delete(ctx, userID, companyID)
}

func (u *Usecase) authorize(ctx context.Context, companyID uuid.UUID, ident identity.Identity) error {
	if ident.Role != identity.RoleHR {
		return identity.ErrHrRoleRequired
	}

	memb, err := u.memberRepo.Get(ctx, ident.UserID, companyID)
	if errors.Is(err, member.ErrCompanyMemberNotFound) {
		return member.ErrCompanyMemberRequired
	}
	if err != nil {
		return err
	}

	if memb.Role != member.CompanyRoleAdmin {
		return member.ErrInsufficientRoleInCompany
	}

	return nil
}
