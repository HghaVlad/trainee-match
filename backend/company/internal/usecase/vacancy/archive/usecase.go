package archive_vacancy

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
)

// Usecase archive vacancy (hide from candidates)
type Usecase struct {
	vacRepo     VacancyRepo
	memberRepo  CompMemberRepo
	compRepo    CompanyRepo
	txManager   uc_common.TxManager
	vacCache    CacheRepo
	pubVacCache CacheRepo
	compCache   CacheRepo
}

func NewUsecase(
	vacRepo VacancyRepo,
	compRepo CompanyRepo,
	memberRepo CompMemberRepo,
	txManager uc_common.TxManager,
	vacCache CacheRepo,
	pubVacCache CacheRepo,
	compCache CacheRepo,
) *Usecase {
	return &Usecase{
		vacRepo:     vacRepo,
		memberRepo:  memberRepo,
		compRepo:    compRepo,
		txManager:   txManager,
		vacCache:    vacCache,
		pubVacCache: pubVacCache,
		compCache:   compCache,
	}
}

// Execute archives vacancy, decreases open vacancies count of company if it was published.
// Deletes company and vacancy from cache.
func (u *Usecase) Execute(
	ctx context.Context,
	compID uuid.UUID,
	vacID uuid.UUID,
	identity uc_common.Identity,
) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		if err := u.authorize(ctx, compID, identity); err != nil {
			return err
		}

		vacancy, err := u.vacRepo.GetByID(ctx, vacID, compID)
		if err != nil {
			return err
		}

		if vacancy.Status == value_types.VacancyStatusArchived {
			return nil
		}

		if err := u.vacRepo.Archive(ctx, compID, vacID); err != nil {
			return err
		}

		if vacancy.Status == value_types.VacancyStatusPublished {
			err := u.compRepo.DecrementOpenVacancies(ctx, vacID)
			if err != nil {
				return err
			}
		}

		u.vacCache.Del(ctx, compID)
		u.pubVacCache.Del(ctx, vacID)
		u.compCache.Del(ctx, vacID)
		return nil
	})
}

// only member of company can archive vacancy
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
