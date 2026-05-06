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

// Usecase archive vacancy (hide from candidates).
// Decreases company OpenVacanciesCount.
// Generates vacancy archived event
type Usecase struct {
	vacRepo      VacancyRepo
	memberRepo   CompMemberRepo
	compRepo     CompanyRepo
	outboxWriter outboxWriter
	txManager    common.TxManager
	vacCache     CacheRepo
	pubVacCache  CacheRepo
	compCache    CacheRepo
}

func NewUsecase(
	vacRepo VacancyRepo,
	compRepo CompanyRepo,
	memberRepo CompMemberRepo,
	outboxWriter outboxWriter,
	txManager common.TxManager,
	vacCache CacheRepo,
	pubVacCache CacheRepo,
	compCache CacheRepo,
) *Usecase {
	return &Usecase{
		vacRepo:      vacRepo,
		memberRepo:   memberRepo,
		compRepo:     compRepo,
		outboxWriter: outboxWriter,
		txManager:    txManager,
		vacCache:     vacCache,
		pubVacCache:  pubVacCache,
		compCache:    compCache,
	}
}

// Execute archives vacancy, decreases open vacancies count of company
// and creates vacancy.ArchivedEvent if it was published.
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

		oldStatus, err := u.vacRepo.ArchiveAndGetOldStatus(ctx, vacID, compID)
		if err != nil {
			return err
		}

		if oldStatus == vacancy.StatusArchived {
			return nil
		}

		if oldStatus == vacancy.StatusPublished {
			err = u.compRepo.DecrementOpenVacancies(ctx, compID)
			if err != nil {
				return err
			}

			err = u.createArchivedEvent(ctx, vacID)
			if err != nil {
				return err
			}

			u.compCache.Del(ctx, compID)
		}

		u.vacCache.Del(ctx, vacID)
		u.pubVacCache.Del(ctx, vacID)
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

func (u *Usecase) createArchivedEvent(ctx context.Context, vacID uuid.UUID) error {
	ev := vacancy.ArchivedEvent{
		EventID:    uuid.New(),
		VacancyID:  vacID,
		OccurredAt: time.Now().UTC(),
	}

	return u.outboxWriter.WriteVacancyArchived(ctx, ev)
}
