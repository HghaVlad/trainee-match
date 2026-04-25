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

func (u *Usecase) Execute(ctx context.Context, id uuid.UUID, identity *identity.Identity) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := u.authorize(ctx, id, identity)
	if err != nil {
		return err
	}

	err = u.compRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	u.cache.Del(ctx, id)
	return nil
}

// only admin of company can delete or admin of the platform
func (u *Usecase) authorize(ctx context.Context, id uuid.UUID, ident *identity.Identity) error {
	if ident.Role == identity.RoleAdmin {
		return nil
	}

	if ident.Role != identity.RoleHR {
		return identity.ErrInsufficientRole
	}

	memb, err := u.memberRepo.Get(ctx, ident.UserID, id)
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
