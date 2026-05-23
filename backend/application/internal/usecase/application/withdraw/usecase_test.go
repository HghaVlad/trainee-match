package withdraw_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/withdraw"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/withdraw/mocks"
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

func newUc(deps *tDeps) *withdraw.Usecase {
	return withdraw.NewUsecase(
		deps.appRepo,
		deps.historyRepo,
		new(fakeTxManager),
	)
}

func candidateIdentity() identity.Identity {
	return identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleCandidate,
	}
}

func TestUsecase_Execute_Success(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUc(deps)

	ident := candidateIdentity()

	appID := uuid.New()

	comment := "found another job"

	req := withdraw.Request{
		AppID:   appID,
		Comment: &comment,
	}

	app := &application.Application{
		ID:        appID,
		Status:    application.StatusSubmitted,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	view := &views.CandidateDetailedView{
		AppID:  appID,
		Status: application.StatusWithdrawn,
	}

	deps.appRepo.EXPECT().
		GetForUpdateByCandidate(gomock.Any(), appID, ident.UserID).
		Return(app, nil)

	deps.appRepo.EXPECT().
		UpdateStatus(
			gomock.Any(),
			appID,
			application.StatusWithdrawn,
			gomock.Any(),
		).
		Return(nil)

	deps.historyRepo.EXPECT().
		Add(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, change application.StatusChange) error {
			require.Equal(t, appID, change.ApplicationID)
			require.Equal(t, application.StatusWithdrawn, change.Status)
			require.Equal(t, application.ActorCandidate, change.ChangedByRole)
			require.Equal(t, &ident.UserID, change.ChangedByUserID)
			require.Equal(t, &comment, change.Comment)

			return nil
		})

	deps.appRepo.EXPECT().
		GetCandidateDetailedView(gomock.Any(), appID, ident.UserID).
		Return(view, nil)

	resp, err := uc.Execute(context.Background(), req, ident)

	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Equal(t, application.StatusWithdrawn, resp.Status)
}

func TestUsecase_Execute_NotCandidate(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUc(deps)

	req := withdraw.Request{
		AppID: uuid.New(),
	}

	ident := identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleHR,
	}

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, identity.ErrCandidateRoleRequired)
	require.Nil(t, resp)
}

func TestUsecase_Execute_InvalidRequest(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUc(deps)

	comment := make([]rune, application.MaxCommentLength+1)
	commentStr := string(comment)

	req := withdraw.Request{
		AppID:   uuid.New(),
		Comment: &commentStr,
	}

	resp, err := uc.Execute(context.Background(), req, candidateIdentity())

	require.ErrorIs(t, err, application.ErrStatusChangeCommentTooLong)
	require.Nil(t, resp)
}

func TestUsecase_Execute_GetForUpdateError(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUc(deps)

	ident := candidateIdentity()

	appID := uuid.New()

	req := withdraw.Request{
		AppID: appID,
	}

	deps.appRepo.EXPECT().
		GetForUpdateByCandidate(gomock.Any(), appID, ident.UserID).
		Return(nil, application.ErrNotFound)

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, application.ErrNotFound)
	require.Nil(t, resp)
}

func TestUsecase_Execute_InvalidStatusTransition(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUc(deps)

	ident := candidateIdentity()

	appID := uuid.New()

	req := withdraw.Request{
		AppID: appID,
	}

	app := &application.Application{
		ID:        appID,
		Status:    application.StatusRejected,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	deps.appRepo.EXPECT().
		GetForUpdateByCandidate(gomock.Any(), appID, ident.UserID).
		Return(app, nil)

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, application.ErrInvalidStatusTransition)
	require.Nil(t, resp)
}

func TestUsecase_Execute_UpdateStatusError(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUc(deps)

	ident := candidateIdentity()

	appID := uuid.New()

	req := withdraw.Request{
		AppID: appID,
	}

	app := &application.Application{
		ID:        appID,
		Status:    application.StatusSubmitted,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	expectedErr := errors.New("update status error")

	deps.appRepo.EXPECT().
		GetForUpdateByCandidate(gomock.Any(), appID, ident.UserID).
		Return(app, nil)

	deps.appRepo.EXPECT().
		UpdateStatus(
			gomock.Any(),
			appID,
			application.StatusWithdrawn,
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

	ident := candidateIdentity()

	appID := uuid.New()

	req := withdraw.Request{
		AppID: appID,
	}

	app := &application.Application{
		ID:        appID,
		Status:    application.StatusSubmitted,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	expectedErr := errors.New("history error")

	deps.appRepo.EXPECT().
		GetForUpdateByCandidate(gomock.Any(), appID, ident.UserID).
		Return(app, nil)

	deps.appRepo.EXPECT().
		UpdateStatus(
			gomock.Any(),
			appID,
			application.StatusWithdrawn,
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

	ident := candidateIdentity()

	appID := uuid.New()

	req := withdraw.Request{
		AppID: appID,
	}

	app := &application.Application{
		ID:        appID,
		Status:    application.StatusSubmitted,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	deps.appRepo.EXPECT().
		GetForUpdateByCandidate(gomock.Any(), appID, ident.UserID).
		Return(app, nil)

	deps.appRepo.EXPECT().
		UpdateStatus(
			gomock.Any(),
			appID,
			application.StatusWithdrawn,
			gomock.Any(),
		).
		Return(nil)

	deps.historyRepo.EXPECT().
		Add(gomock.Any(), gomock.Any()).
		Return(nil)

	deps.appRepo.EXPECT().
		GetCandidateDetailedView(gomock.Any(), appID, ident.UserID).
		Return(nil, application.ErrNotFound)

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, application.ErrNotFound)
	require.Nil(t, resp)
}
