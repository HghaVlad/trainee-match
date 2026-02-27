package delete_company

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

func (u *Usecase) Execute(ctx context.Context, id uuid.UUID, identity uc_common.Identity) error {
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
func (u *Usecase) authorize(ctx context.Context, id uuid.UUID, identity uc_common.Identity) error {
	if identity.Role == uc_common.RoleAdmin {
		return nil
	}

	if identity.Role != uc_common.RoleHR {
		return domain_errors.ErrInsufficientRole
	}

	member, err := u.memberRepo.Get(ctx, identity.UserID, id)
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
