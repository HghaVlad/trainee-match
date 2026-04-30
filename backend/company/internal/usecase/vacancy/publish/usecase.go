package publish

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

// Usecase publishes vacancy, thus makes it available for candidates.
// Increases company open vacancies count
// Generates a vacancy.PublishedEvent if it isn't in published state now
type Usecase struct {
	vacRepo      VacancyRepo
	companyRepo  CompanyRepo
	memberRepo   CompMemberRepo
	outboxWriter outboxWriter
	txManager    common.TxManager
	compCache    CacheRepo
	vacCache     CacheRepo
}

func NewUsecase(
	vacRepo VacancyRepo,
	compRepo CompanyRepo,
	memberRepo CompMemberRepo,
	outboxWriter outboxWriter,
	txManager common.TxManager,
	vacCache CacheRepo,
	compCache CacheRepo,
) *Usecase {
	return &Usecase{
		vacRepo:      vacRepo,
		memberRepo:   memberRepo,
		companyRepo:  compRepo,
		outboxWriter: outboxWriter,
		txManager:    txManager,
		compCache:    compCache,
		vacCache:     vacCache,
	}
}

// Execute publishes vacancy, thus makes it available for candidates.
// Increases company open vacancies count and creates vacancy.PublishedEvent, if vacancy wasn't published.
// Returns vacacny.ErrNotFound if it was.
// Deletes company and vacancy from cache because of the updates.
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

		vac, err := u.vacRepo.PublishIfNotPublished(ctx, vacID, compID)
		if err != nil {
			return err
		}

		if vac.WasAlreadyPublished {
			return nil
		}

		err = u.companyRepo.IncrementOpenVacancies(ctx, compID)
		if err != nil {
			return err
		}

		err = u.createPublishedEvent(ctx, vac)
		if err != nil {
			return err
		}

		u.compCache.Del(ctx, compID)
		u.vacCache.Del(ctx, vacID)
		return nil
	})
}

// only member of company can publish vacancy
func (u *Usecase) authorize(ctx context.Context, companyID uuid.UUID, ident *identity.Identity) error {
	if ident.Role != identity.RoleHR {
		return identity.ErrHrRoleRequired
	}

	_, err := u.memberRepo.Get(ctx, ident.UserID, companyID)
	if errors.Is(err, member.ErrCompanyMemberNotFound) {
		return member.ErrCompanyMemberRequired
	}

	return err
}

func (u *Usecase) createPublishedEvent(ctx context.Context, vac *PublishedEventView) error {
	ev := vacancy.PublishedEvent{
		EventID:     uuid.New(),
		VacancyID:   vac.ID,
		Title:       vac.Title,
		CompanyID:   vac.CompanyID,
		CompanyName: vac.CompanyName,
		OccurredAt:  time.Now().UTC(),
	}

	return u.outboxWriter.WriteVacancyPublished(ctx, ev)
}
