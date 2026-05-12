package hrupdstatus

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
) (*views.HrDetailedView, error) {
	if ident.Role != identity.RoleHR {
		return nil, identity.ErrHrRoleRequired
	}

	if err := req.validate(); err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	err := u.txManager.Do(ctx, func(ctx context.Context) error {
		app, err := u.appRepo.GetForUpdateByHr(ctx, req.AppID, ident.UserID)
		if err != nil {
			return err
		}

		err = app.ChangeStatus(req.Status, application.ActorHR, now)
		if err != nil {
			return err
		}

		statusChange, err := application.NewStatusChange(
			app.ID,
			app.Status,
			&ident.UserID,
			application.ActorHR,
			req.Comment,
			now,
		)
		if err != nil {
			return err
		}

		err = u.appRepo.UpdateStatus(ctx, app.ID, app.Status, now)
		if err != nil {
			return err
		}

		return u.appStatusHistoryRepo.Add(ctx, *statusChange)
	})

	if err != nil {
		return nil, err
	}

	view, err := u.appRepo.GetHrDetailedView(ctx, req.AppID, ident.UserID)
	if err != nil {
		return nil, err
	}

	view.AllowedActions = views.HRAllowedActions(view.Status)
	return view, nil
}
