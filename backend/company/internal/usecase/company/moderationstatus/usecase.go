package moderationstatus

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
)

type Usecase struct {
	compRepo     companyRepo
	outboxWriter outboxWriter
	txManager    common.TxManager
	cache        cacheRepo
}

func NewUsecase(
	repo companyRepo,
	outboxWriter outboxWriter,
	txManager common.TxManager,
	cache cacheRepo,
) *Usecase {
	return &Usecase{
		compRepo:     repo,
		outboxWriter: outboxWriter,
		txManager:    txManager,
		cache:        cache,
	}
}

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

	now := time.Now().UTC()
	upd := false

	if err := u.authorize(identity); err != nil {
		return err
	}

	err := u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		oldStatus, err := u.compRepo.UpdateModerationStatusAndGetOld(ctx, req.ID, req.Status, now)
		if err != nil {
			return err
		}

		if oldStatus == req.Status {
			return nil
		}

		err = u.createModerationUpdatedEvent(ctx, req.ID, req.Status, now)
		if err != nil {
			return err
		}

		upd = true
		return nil
	})

	if err != nil {
		return err
	}

	if upd {
		u.cache.Del(ctx, req.ID)
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
	companyID uuid.UUID,
	status company.ModerationStatus,
	when time.Time,
) error {
	ev := company.ModerationUpdatedEvent{
		EventID:          uuid.New(),
		CompanyID:        companyID,
		ModerationStatus: status,
		OccurredAt:       when,
	}

	return u.outboxWriter.WriteCompanyModerationUpdated(ctx, ev)
}
