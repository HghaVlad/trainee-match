package get

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
)

// Usecase returns full model of vacancy (for company members)
type Usecase struct {
	repo       Repository
	cache      CacheRepo
	memberRepo CompMemberRepo
}

func NewUsecase(repo Repository, cache CacheRepo, memberRepo CompMemberRepo) *Usecase {
	return &Usecase{
		repo:       repo,
		cache:      cache,
		memberRepo: memberRepo,
	}
}

// Execute returns full model of vacancy (for company members)
func (u *Usecase) Execute(
	ctx context.Context,
	vacancyID uuid.UUID,
	companyID uuid.UUID,
	ident identity.Identity,
) (*vacancy.Vacancy, error) {
	if err := u.authorize(ctx, companyID, ident); err != nil {
		return nil, err
	}

	vac := u.cache.Get(ctx, vacancyID)

	if vac != nil && vac.CompanyID == companyID {
		return vac, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	vac, err := u.repo.GetByID(ctx, vacancyID, companyID)
	if err != nil {
		return nil, err
	}

	u.cache.Put(ctx, vacancyID, vac, time.Second*300)
	return vac, nil
}

// only member of company has access
func (u *Usecase) authorize(ctx context.Context, companyID uuid.UUID, ident identity.Identity) error {
	if ident.Role != identity.RoleHR {
		return identity.ErrHrRoleRequired
	}

	_, err := u.memberRepo.Get(ctx, ident.UserID, companyID)
	if errors.Is(err, member.ErrCompanyMemberNotFound) {
		return member.ErrCompanyMemberRequired
	}

	return err
}
