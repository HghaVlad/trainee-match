package moderationstatus_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/moderationstatus"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/moderationstatus/mocks"
)

type fakeTxManager struct{}

func (f *fakeTxManager) WithinTx(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	return fn(ctx)
}

func TestUsecase_Execute_Success_StatusChanged(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	compRepo := mocks.NewMockcompanyRepo(ctrl)
	outboxWriter := mocks.NewMockoutboxWriter(ctrl)
	cache := mocks.NewMockcacheRepo(ctrl)

	uc := moderationstatus.NewUsecase(
		compRepo,
		outboxWriter,
		&fakeTxManager{},
		cache,
	)

	req := &moderationstatus.Request{
		ID:     uuid.New(),
		Status: company.ModerationStatusOK,
	}

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleAdmin,
	}

	compRepo.
		EXPECT().
		UpdateModerationStatusAndGetOld(
			gomock.Any(),
			req.ID,
			req.Status,
			gomock.Any(),
		).
		Return(company.ModerationStatusHidden, nil)

	outboxWriter.
		EXPECT().
		WriteCompanyModerationUpdated(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, ev company.ModerationUpdatedEvent) error {
			require.Equal(t, req.ID, ev.CompanyID)
			require.Equal(t, req.Status, ev.ModerationStatus)
			require.NotEqual(t, uuid.Nil, ev.EventID)

			return nil
		})

	cache.
		EXPECT().
		Del(gomock.Any(), req.ID)

	err := uc.Execute(context.Background(), req, ident)

	require.NoError(t, err)
}

func TestUsecase_Execute_Success_StatusNotChanged(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	compRepo := mocks.NewMockcompanyRepo(ctrl)
	outboxWriter := mocks.NewMockoutboxWriter(ctrl)
	cache := mocks.NewMockcacheRepo(ctrl)

	uc := moderationstatus.NewUsecase(
		compRepo,
		outboxWriter,
		&fakeTxManager{},
		cache,
	)

	req := &moderationstatus.Request{
		ID:     uuid.New(),
		Status: company.ModerationStatusOK,
	}

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleAdmin,
	}

	compRepo.
		EXPECT().
		UpdateModerationStatusAndGetOld(
			gomock.Any(),
			req.ID,
			req.Status,
			gomock.Any(),
		).
		Return(company.ModerationStatusOK, nil)

	err := uc.Execute(context.Background(), req, ident)

	require.NoError(t, err)
}

func TestUsecase_Execute_InvalidStatus(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	uc := moderationstatus.NewUsecase(
		mocks.NewMockcompanyRepo(ctrl),
		mocks.NewMockoutboxWriter(ctrl),
		&fakeTxManager{},
		mocks.NewMockcacheRepo(ctrl),
	)

	req := &moderationstatus.Request{
		ID:     uuid.New(),
		Status: company.ModerationStatus("invalid"),
	}

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleAdmin,
	}

	err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, company.ErrInvalidModerationStatus)
}

func TestUsecase_Execute_AdminRoleRequired(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	uc := moderationstatus.NewUsecase(
		mocks.NewMockcompanyRepo(ctrl),
		mocks.NewMockoutboxWriter(ctrl),
		&fakeTxManager{},
		mocks.NewMockcacheRepo(ctrl),
	)

	req := &moderationstatus.Request{
		ID:     uuid.New(),
		Status: company.ModerationStatusOK,
	}

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleHR,
	}

	err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, identity.ErrAdminRoleRequired)
}

func TestUsecase_Execute_UpdateModStatusError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	compRepo := mocks.NewMockcompanyRepo(ctrl)

	uc := moderationstatus.NewUsecase(
		compRepo,
		mocks.NewMockoutboxWriter(ctrl),
		&fakeTxManager{},
		mocks.NewMockcacheRepo(ctrl),
	)

	req := &moderationstatus.Request{
		ID:     uuid.New(),
		Status: company.ModerationStatusOK,
	}

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleAdmin,
	}

	compRepo.
		EXPECT().
		UpdateModerationStatusAndGetOld(
			gomock.Any(),
			req.ID,
			req.Status,
			gomock.Any(),
		).
		Return(company.ModerationStatusHidden, company.ErrCompanyNotFound)

	err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, company.ErrCompanyNotFound)
}

func TestUsecase_Execute_OutboxError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	compRepo := mocks.NewMockcompanyRepo(ctrl)
	outboxWriter := mocks.NewMockoutboxWriter(ctrl)

	uc := moderationstatus.NewUsecase(
		compRepo,
		outboxWriter,
		&fakeTxManager{},
		mocks.NewMockcacheRepo(ctrl),
	)

	req := &moderationstatus.Request{
		ID:     uuid.New(),
		Status: company.ModerationStatusOK,
	}

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleAdmin,
	}

	expectedErr := errors.New("outbox failed")

	compRepo.
		EXPECT().
		UpdateModerationStatusAndGetOld(
			gomock.Any(),
			req.ID,
			req.Status,
			gomock.Any(),
		).
		Return(company.ModerationStatusHidden, nil)

	outboxWriter.
		EXPECT().
		WriteCompanyModerationUpdated(gomock.Any(), gomock.Any()).
		Return(expectedErr)

	err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, expectedErr)
}

func TestUsecase_Execute_CacheDeletedOnlyWhenUpdated(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	compRepo := mocks.NewMockcompanyRepo(ctrl)
	outboxWriter := mocks.NewMockoutboxWriter(ctrl)
	cache := mocks.NewMockcacheRepo(ctrl)

	uc := moderationstatus.NewUsecase(
		compRepo,
		outboxWriter,
		&fakeTxManager{},
		cache,
	)

	req := &moderationstatus.Request{
		ID:     uuid.New(),
		Status: company.ModerationStatusHidden,
	}

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleAdmin,
	}

	compRepo.
		EXPECT().
		UpdateModerationStatusAndGetOld(
			gomock.Any(),
			req.ID,
			req.Status,
			gomock.Any(),
		).
		Return(company.ModerationStatusHidden, nil)

	err := uc.Execute(context.Background(), req, ident)

	require.NoError(t, err)
}
