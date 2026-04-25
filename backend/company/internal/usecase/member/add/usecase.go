package add

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

func (u *Usecase) Execute(ctx context.Context, req *Request, identity *identity.Identity) error {
	if err := req.Validate(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := u.authorize(ctx, req.CompanyID, identity); err != nil {
		return err
	}

	memb := &member.CompanyMember{
		UserID:    req.UserID,
		CompanyID: req.CompanyID,
		Role:      req.Role,
	}

	return u.memberRepo.Create(ctx, memb)
}

func (u *Usecase) authorize(ctx context.Context, companyID uuid.UUID, ident *identity.Identity) error {
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
