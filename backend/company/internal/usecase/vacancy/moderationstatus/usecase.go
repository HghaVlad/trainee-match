package moderationstatus

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
)

type Usecase struct {
	vacRepo      vacancyRepo
	compRepo     companyRepo
	outboxWriter outboxWriter
	txManager    common.TxManager
	searchRepo   searchRepo
	vacCache     cacheRepo
	pubVacCache  cacheRepo
	compCache    cacheRepo
}

func NewUsecase(
	vacRepo vacancyRepo,
	compRepo companyRepo,
	outboxWriter outboxWriter,
	txManager common.TxManager,
	searchRepo searchRepo,
	vacCache cacheRepo,
	pubVacCache cacheRepo,
	compCache cacheRepo,
) *Usecase {
	return &Usecase{
		vacRepo:      vacRepo,
		compRepo:     compRepo,
		outboxWriter: outboxWriter,
		txManager:    txManager,
		searchRepo:   searchRepo,
		vacCache:     vacCache,
		pubVacCache:  pubVacCache,
		compCache:    compCache,
	}
}

// Execute updates moderation status, publishes event if it was different previously.
// Works only if vacancy is published, returns vacancy.ErrVacancyNotFound otherwise.
// Also keeps company open vacancies count up to date.
func (u *Usecase) Execute(
	ctx context.Context,
	req *Request,
	identity *identity.Identity,
) error {
	if err := req.Validate(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	if err := u.authorize(identity); err != nil {
		return err
	}

	now := time.Now().UTC()
	newModStatus := req.Status
	vacUpd := false
	compUpd := false
	var compID uuid.UUID

	err := u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		res, err := u.vacRepo.UpdateModerationStatus(ctx, req.ID, newModStatus, now)
		if err != nil {
			return err
		}

		if res.OldModerationStatus == newModStatus {
			return nil
		}

		compUpd, err = u.updateCompanyOpenVacCount(ctx, *res, newModStatus)
		if err != nil {
			return err
		}
		compID = res.CompanyID

		err = u.createModerationUpdatedEvent(ctx, req.ID, newModStatus, now)
		if err != nil {
			return err
		}

		vacUpd = true
		return nil
	})

	if err != nil {
		return err
	}

	if vacUpd {
		u.vacCache.Del(ctx, req.ID)
		u.pubVacCache.Del(ctx, req.ID)

		searchView, err := u.vacRepo.GetSearchView(ctx, req.ID)
		if err != nil {
			return err
		}

		if err := u.searchRepo.Index(ctx, *searchView); err != nil {
			return err
		}
	}

	if compUpd {
		u.compCache.Del(ctx, compID)
	}

	return nil
}

func (u *Usecase) authorize(ident *identity.Identity) error {
	if ident.Role != identity.RoleAdmin {
		return identity.ErrAdminRoleRequired
	}

	return nil
}

func (u *Usecase) createModerationUpdatedEvent(
	ctx context.Context,
	vacID uuid.UUID,
	status vacancy.ModerationStatus,
	when time.Time,
) error {
	ev := vacancy.ModerationUpdatedEvent{
		EventID:          uuid.New(),
		VacancyID:        vacID,
		ModerationStatus: status,
		OccurredAt:       when,
	}

	return u.outboxWriter.WriteVacancyModerationUpdated(ctx, ev)
}

func (u *Usecase) updateCompanyOpenVacCount(
	ctx context.Context,
	res UpdateModerationResult,
	newModStatus vacancy.ModerationStatus,
) (bool, error) {
	if res.VacancyStatus == vacancy.StatusPublished &&
		newModStatus == vacancy.ModerationStatusHidden {
		err := u.compRepo.DecrementOpenVacancies(ctx, res.CompanyID)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	if res.VacancyStatus == vacancy.StatusPublished &&
		newModStatus == vacancy.ModerationStatusOK {
		err := u.compRepo.IncrementOpenVacancies(ctx, res.CompanyID)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}
