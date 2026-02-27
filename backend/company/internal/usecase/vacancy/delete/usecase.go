package delete_vacancy

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
)

type Usecase struct {
	compRepo   VacancyRepo
	memberRepo CompMemberRepo
	cache      CacheRepo
}

func NewUsecase(repo VacancyRepo, memberRepo CompMemberRepo, cache CacheRepo) *Usecase {
	return &Usecase{compRepo: repo, memberRepo: memberRepo, cache: cache}
}

func (u *Usecase) Execute(
	ctx context.Context,
	vacancyID uuid.UUID,
	companyID uuid.UUID,
	identity uc_common.Identity,
) error {

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := u.authorize(ctx, companyID, identity); err != nil {
		return err
	}

	err := u.compRepo.Delete(ctx, vacancyID, companyID)
	if err != nil {
		return err
	}

	u.cache.Del(ctx, vacancyID)
	return nil
}

// only member of company can delete vacancy
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
