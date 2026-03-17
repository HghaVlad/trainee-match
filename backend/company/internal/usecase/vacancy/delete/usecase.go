package delete_vacancy

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
)

// Usecase hard delete from db. Reckon using archive vacancy instead.
type Usecase struct {
	vacRepo     VacancyRepo
	compRepo    CompanyRepo
	txManager   uc_common.TxManager
	memberRepo  CompMemberRepo
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

// Execute hard delete from db. Reckon using archive vacancy instead.
// Also removes company from cache because of the updated OpenVacCnt
func (u *Usecase) Execute(
	ctx context.Context,
	vacancyID uuid.UUID,
	companyID uuid.UUID,
	identity uc_common.Identity,
) error {

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		if err := u.authorize(ctx, companyID, identity); err != nil {
			return err
		}

		vacancy, err := u.vacRepo.GetByID(ctx, vacancyID, companyID)
		if err != nil {
			return err
		}

		if err := u.vacRepo.Delete(ctx, vacancyID, companyID); err != nil {
			return err
		}

		if vacancy.Status == value_types.VacancyStatusPublished {
			if err := u.compRepo.DecrementOpenVacancies(ctx, companyID); err != nil {
				return err
			}
		}

		u.vacCache.Del(ctx, vacancyID)
		u.pubVacCache.Del(ctx, vacancyID)
		u.compCache.Del(ctx, companyID)
		return nil
	})
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
