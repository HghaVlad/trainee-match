package update

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
)

type Usecase struct {
	compRepo   CompanyRepo
	memberRepo CompMemberRepo
	cache      CacheRepo
}

func NewUsecase(
	repo CompanyRepo,
	memberRepo CompMemberRepo,
	cache CacheRepo,
) *Usecase {
	return &Usecase{
		compRepo:   repo,
		memberRepo: memberRepo,
		cache:      cache,
	}
}

func (u *Usecase) Execute(ctx context.Context, req *Request, identity *identity.Identity) error {
	if err := req.Validate(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := u.authorize(ctx, req.ID, identity)
	if err != nil {
		return err
	}

	err = u.compRepo.Update(ctx, req)
	if err != nil {
		return err
	}

	u.cache.Del(ctx, req.ID)
	return nil
}

// only admin of company can update
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
