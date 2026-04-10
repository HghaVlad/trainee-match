package archive

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
)

// Usecase archive vacancy (hide from candidates)
type Usecase struct {
	vacRepo     VacancyRepo
	memberRepo  CompMemberRepo
	compRepo    CompanyRepo
	txManager   common.TxManager
	vacCache    CacheRepo
	pubVacCache CacheRepo
	compCache   CacheRepo
}

func NewUsecase(
	vacRepo VacancyRepo,
	compRepo CompanyRepo,
	memberRepo CompMemberRepo,
	txManager common.TxManager,
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
	identity *identity.Identity,
) error {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	return u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		if err := u.authorize(ctx, compID, identity); err != nil {
			return err
		}

		vac, err := u.vacRepo.GetByID(ctx, vacID, compID)
		if err != nil {
			return err
		}

		if vac.Status == vacancy.StatusArchived {
			return nil
		}

		if err := u.vacRepo.Archive(ctx, vacID, compID); err != nil {
			return err
		}

		if vac.Status == vacancy.StatusPublished {
			err := u.compRepo.DecrementOpenVacancies(ctx, compID)
			if err != nil {
				return err
			}
		}

		u.vacCache.Del(ctx, vacID)
		u.pubVacCache.Del(ctx, vacID)
		u.compCache.Del(ctx, compID)
		return nil
	})
}

// only member of company can archive vacancy
func (u *Usecase) authorize(ctx context.Context, companyID uuid.UUID, iden *identity.Identity) error {
	if iden.Role != identity.RoleHR {
		return identity.ErrHrRoleRequired
	}

	_, err := u.memberRepo.Get(ctx, iden.UserID, companyID)
	if errors.Is(err, member.ErrCompanyMemberNotFound) {
		return member.ErrCompanyMemberRequired
	}

	return err
}
