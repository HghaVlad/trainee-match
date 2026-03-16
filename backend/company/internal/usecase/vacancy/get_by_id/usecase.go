package get_vacancy

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
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
	identity uc_common.Identity,
) (*domain.Vacancy, error) {

	if err := u.authorize(ctx, companyID, identity); err != nil {
		return nil, err
	}

	vacancy := u.cache.Get(ctx, vacancyID)

	if vacancy != nil && vacancy.CompanyID == companyID {
		return vacancy, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	vacancy, err := u.repo.GetByID(ctx, vacancyID, companyID)
	if err != nil {
		return nil, err
	}

	u.cache.Put(ctx, vacancyID, vacancy, time.Second*300)
	return vacancy, nil
}

// only member of company has access
func (u *Usecase) authorize(ctx context.Context, companyID uuid.UUID, identity uc_common.Identity) error {
	if identity.Role != uc_common.RoleHR {
		return domain_errors.ErrHrRoleRequired
	}

	_, err := u.memberRepo.Get(ctx, identity.UserID, companyID)
	if errors.Is(err, domain_errors.ErrCompanyMemberNotFound) {
		return domain_errors.ErrCompanyMemberRequired
	}

	return err
}
