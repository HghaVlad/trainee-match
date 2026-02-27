package update_company

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

func (u *Usecase) Execute(ctx context.Context, req *Request, identity uc_common.Identity) error {
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
