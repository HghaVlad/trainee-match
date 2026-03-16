package publish_vacancy

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
)

// Usecase publishes vacancy, thus makes it available for candidates.
// Increases company open vacancies count
type Usecase struct {
	vacRepo     VacancyRepo
	companyRepo CompanyRepo
	memberRepo  CompMemberRepo
	txManager   uc_common.TxManager
	compCache   CacheRepo
	vacCache    CacheRepo
}

func NewUsecase(
	vacRepo VacancyRepo,
	compRepo CompanyRepo,
	memberRepo CompMemberRepo,
	txManager uc_common.TxManager,
	vacCache CacheRepo,
	compCache CacheRepo,
) *Usecase {

	return &Usecase{
		vacRepo:     vacRepo,
		memberRepo:  memberRepo,
		companyRepo: compRepo,
		txManager:   txManager,
		compCache:   compCache,
		vacCache:    vacCache,
	}
}

// Execute publishes vacancy, thus makes it available for candidates.
// Increases company open vacancies count if vacancy wasn't published.
// Deletes company and vacancy from cache because of the updates.
func (u *Usecase) Execute(
	ctx context.Context,
	compID uuid.UUID,
	vacID uuid.UUID,
	identity uc_common.Identity,
) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	return u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		if err := u.authorize(ctx, compID, identity); err != nil {
			return err
		}

		vacancy, err := u.vacRepo.GetByID(ctx, vacID, compID)
		if err != nil {
			return err
		}

		if vacancy.Status == value_types.VacancyStatusPublished {
			return nil
		}

		err = u.vacRepo.Publish(ctx, compID, vacID)
		if err != nil {
			return err
		}

		err = u.companyRepo.IncrementOpenVacancies(ctx, compID)
		if err != nil {
			return err
		}

		u.compCache.Del(ctx, compID)
		u.vacCache.Del(ctx, vacID)
		return nil
	})
}

// only member of company can publish vacancy
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
