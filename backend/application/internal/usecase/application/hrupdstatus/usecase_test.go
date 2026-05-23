package hrupdstatus_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/hrupdstatus"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/hrupdstatus/mocks"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/identity"
)

type fakeTxManager struct{}

func (f *fakeTxManager) Do(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	return fn(ctx)
}

type tDeps struct {
	appRepo     *mocks.MockappRepo
	historyRepo *mocks.MockappStatusHistoryRepo
}

func setup(t *testing.T) *tDeps {
	ctrl := gomock.NewController(t)

	return &tDeps{
		appRepo:     mocks.NewMockappRepo(ctrl),
		historyRepo: mocks.NewMockappStatusHistoryRepo(ctrl),
	}
}

func newUc(deps *tDeps) *hrupdstatus.Usecase {
	return hrupdstatus.NewUsecase(
		deps.appRepo,
		deps.historyRepo,
		new(fakeTxManager),
	)
}

func hrIdentity() identity.Identity {
	return identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleHR,
	}
}

func TestUsecase_Execute_Success(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUc(deps)

	ident := hrIdentity()
	appID := uuid.New()
	comment := "looks good"

	req := hrupdstatus.Request{
		AppID:   appID,
		Status:  application.StatusInterview,
		Comment: &comment,
	}

	app := &application.Application{
		ID:        appID,
		Status:    application.StatusSubmitted,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	view := &views.HrDetailedView{
		AppID:  appID,
		Status: application.StatusInterview,
	}

	deps.appRepo.EXPECT().
		GetForUpdateByHr(gomock.Any(), appID, ident.UserID).
		Return(app, nil)

	deps.appRepo.EXPECT().
		UpdateStatus(
			gomock.Any(),
			appID,
			application.StatusInterview,
			gomock.Any(),
		).
		Return(nil)

	deps.historyRepo.EXPECT().
		Add(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, change application.StatusChange) error {
			require.Equal(t, appID, change.ApplicationID)
			require.Equal(t, application.StatusInterview, change.Status)
			require.Equal(t, application.ActorHR, change.ChangedByRole)
			require.Equal(t, &ident.UserID, change.ChangedByUserID)
			require.Equal(t, &comment, change.Comment)

			return nil
		})

	deps.appRepo.EXPECT().
		GetHrDetailedView(gomock.Any(), appID, ident.UserID).
		Return(view, nil)

	resp, err := uc.Execute(context.Background(), req, ident)

	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Equal(t, application.StatusInterview, resp.Status)
	require.NotNil(t, resp.AllowedActions)
}

func TestUsecase_Execute_NotHR(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUc(deps)

	req := hrupdstatus.Request{
		AppID:  uuid.New(),
		Status: application.StatusInterview,
	}

	ident := identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleCandidate,
	}

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, identity.ErrHrRoleRequired)
	require.Nil(t, resp)
}

func TestUsecase_Execute_InvalidRequest(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUc(deps)

	req := hrupdstatus.Request{
		AppID:  uuid.New(),
		Status: application.Status("invalid"),
	}

	resp, err := uc.Execute(context.Background(), req, hrIdentity())

	require.ErrorIs(t, err, application.ErrInvalidStatus)
	require.Nil(t, resp)
}

func TestUsecase_Execute_GetForUpdateError(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUc(deps)

	ident := hrIdentity()

	appID := uuid.New()

	req := hrupdstatus.Request{
		AppID:  appID,
		Status: application.StatusInterview,
	}

	deps.appRepo.EXPECT().
		GetForUpdateByHr(gomock.Any(), appID, ident.UserID).
		Return(nil, projection.ErrVacancyNotFound)

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, projection.ErrVacancyNotFound)
	require.Nil(t, resp)
}

func TestUsecase_Execute_InvalidStatusTransition(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUc(deps)

	ident := hrIdentity()

	appID := uuid.New()

	req := hrupdstatus.Request{
		AppID:  appID,
		Status: application.StatusSeen,
	}

	app := &application.Application{
		ID:        appID,
		Status:    application.StatusInterview,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	deps.appRepo.EXPECT().
		GetForUpdateByHr(gomock.Any(), appID, ident.UserID).
		Return(app, nil)

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, application.ErrInvalidStatusTransition)
	require.Nil(t, resp)
}

func TestUsecase_Execute_UpdateStatusError(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUc(deps)

	ident := hrIdentity()

	appID := uuid.New()

	req := hrupdstatus.Request{
		AppID:  appID,
		Status: application.StatusInterview,
	}

	app := &application.Application{
		ID:        appID,
		Status:    application.StatusSubmitted,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	expectedErr := errors.New("update status error")

	deps.appRepo.EXPECT().
		GetForUpdateByHr(gomock.Any(), appID, ident.UserID).
		Return(app, nil)

	deps.appRepo.EXPECT().
		UpdateStatus(
			gomock.Any(),
			appID,
			application.StatusInterview,
			gomock.Any(),
		).
		Return(expectedErr)

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, resp)
}

func TestUsecase_Execute_AddHistoryError(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUc(deps)

	ident := hrIdentity()

	appID := uuid.New()

	req := hrupdstatus.Request{
		AppID:  appID,
		Status: application.StatusInterview,
	}

	app := &application.Application{
		ID:        appID,
		Status:    application.StatusSubmitted,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	expectedErr := errors.New("history error")

	deps.appRepo.EXPECT().
		GetForUpdateByHr(gomock.Any(), appID, ident.UserID).
		Return(app, nil)

	deps.appRepo.EXPECT().
		UpdateStatus(
			gomock.Any(),
			appID,
			application.StatusInterview,
			gomock.Any(),
		).
		Return(nil)

	deps.historyRepo.EXPECT().
		Add(gomock.Any(), gomock.Any()).
		Return(expectedErr)

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, resp)
}

func TestUsecase_Execute_GetViewError(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUc(deps)

	ident := hrIdentity()

	appID := uuid.New()

	req := hrupdstatus.Request{
		AppID:  appID,
		Status: application.StatusInterview,
	}

	app := &application.Application{
		ID:        appID,
		Status:    application.StatusSubmitted,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	deps.appRepo.EXPECT().
		GetForUpdateByHr(gomock.Any(), appID, ident.UserID).
		Return(app, nil)

	deps.appRepo.EXPECT().
		UpdateStatus(
			gomock.Any(),
			appID,
			application.StatusInterview,
			gomock.Any(),
		).
		Return(nil)

	deps.historyRepo.EXPECT().
		Add(gomock.Any(), gomock.Any()).
		Return(nil)

	deps.appRepo.EXPECT().
		GetHrDetailedView(gomock.Any(), appID, ident.UserID).
		Return(nil, projection.ErrVacancyNotFound)

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, projection.ErrVacancyNotFound)
	require.Nil(t, resp)
}
