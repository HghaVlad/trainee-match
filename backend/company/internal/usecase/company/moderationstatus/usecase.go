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

	return u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		err := u.authorize(identity)
		if err != nil {
			return err
		}

		err = u.compRepo.UpdateModerationStatus(ctx, req.ID, req.Status)
		if err != nil {
			return err
		}

		err = u.createModerationUpdatedEvent(
			ctx,
			req.ID,
			req.Status,
		)
		if err != nil {
			return err
		}

		u.cache.Del(ctx, req.ID)

		return nil
	})
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
) error {
	ev := company.ModerationUpdatedEvent{
		EventID:          uuid.New(),
		CompanyID:        companyID,
		ModerationStatus: status,
		OccurredAt:       time.Now().UTC(),
	}

	return u.outboxWriter.WriteCompanyModerationUpdated(ctx, ev)
}
