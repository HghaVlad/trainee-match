package withdraw

import (
	"context"
	"time"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/identity"
)

type Usecase struct {
	appRepo              appRepo
	appStatusHistoryRepo appStatusHistoryRepo
	txManager            common.TxManager
}

func NewUsecase(
	appRepo appRepo,
	appStatusHistoryRepo appStatusHistoryRepo,
	txManager common.TxManager,
) *Usecase {
	return &Usecase{
		appRepo:              appRepo,
		appStatusHistoryRepo: appStatusHistoryRepo,
		txManager:            txManager,
	}
}

func (u *Usecase) Execute(
	ctx context.Context,
	req Request,
	ident identity.Identity,
) (*views.CandidateDetailedView, error) {
	if ident.Role != identity.RoleCandidate {
		return nil, identity.ErrCandidateRoleRequired
	}

	if err := req.validate(); err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	err := u.txManager.Do(ctx, func(ctx context.Context) error {
		app, err := u.appRepo.GetForUpdateByCandidate(ctx, req.AppID, ident.UserID)
		if err != nil {
			return err
		}

		err = app.Withdraw(now)
		if err != nil {
			return err
		}

		statusChange, err := application.NewStatusChange(
			req.AppID,
			application.StatusWithdrawn,
			&ident.UserID,
			application.ActorCandidate,
			req.Comment,
			now,
		)
		if err != nil {
			return err
		}

		err = u.appRepo.UpdateStatus(ctx, req.AppID, application.StatusWithdrawn, now)
		if err != nil {
			return err
		}

		return u.appStatusHistoryRepo.Add(ctx, *statusChange)
	})

	if err != nil {
		return nil, err
	}

	return u.appRepo.GetCandidateDetailedView(ctx, req.AppID, ident.UserID)
}
