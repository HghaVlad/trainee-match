package moderationstatus_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/moderationstatus"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/moderationstatus/mocks"
)

type fakeTxManager struct{}

func (f *fakeTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestUsecase_Execute_Success_StatusChanged_Hidden(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	vacRepo := mocks.NewMockvacancyRepo(ctrl)
	compRepo := mocks.NewMockcompanyRepo(ctrl)
	outboxWriter := mocks.NewMockoutboxWriter(ctrl)

	vacCache := mocks.NewMockcacheRepo(ctrl)
	pubVacCache := mocks.NewMockcacheRepo(ctrl)
	compCache := mocks.NewMockcacheRepo(ctrl)

	uc := moderationstatus.NewUsecase(
		vacRepo,
		compRepo,
		outboxWriter,
		&fakeTxManager{},
		vacCache,
		pubVacCache,
		compCache,
	)

	req := &moderationstatus.Request{
		ID:     uuid.New(),
		Status: vacancy.ModerationStatusHidden,
	}

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleAdmin,
	}

	compID := uuid.New()

	vacRepo.
		EXPECT().
		UpdateModerationStatus(
			gomock.Any(),
			req.ID,
			req.Status,
			gomock.Any(),
		).
		Return(&moderationstatus.UpdateModerationResult{
			CompanyID:           compID,
			VacancyStatus:       vacancy.StatusPublished,
			OldModerationStatus: vacancy.ModerationStatusOK,
		}, nil)

	compRepo.
		EXPECT().
		DecrementOpenVacancies(gomock.Any(), compID).
		Return(nil)

	outboxWriter.
		EXPECT().
		WriteVacancyModerationUpdated(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, ev vacancy.ModerationUpdatedEvent) error {
			require.Equal(t, req.ID, ev.VacancyID)
			require.Equal(t, req.Status, ev.ModerationStatus)
			require.NotEqual(t, uuid.Nil, ev.EventID)

			return nil
		})

	vacCache.EXPECT().Del(gomock.Any(), req.ID)

	pubVacCache.EXPECT().Del(gomock.Any(), req.ID)

	compCache.EXPECT().Del(gomock.Any(), compID)

	err := uc.Execute(context.Background(), req, ident)

	require.NoError(t, err)
}

func TestUsecase_Execute_Success_StatusChanged_OK(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	vacRepo := mocks.NewMockvacancyRepo(ctrl)
	compRepo := mocks.NewMockcompanyRepo(ctrl)
	outboxWriter := mocks.NewMockoutboxWriter(ctrl)

	vacCache := mocks.NewMockcacheRepo(ctrl)
	pubVacCache := mocks.NewMockcacheRepo(ctrl)
	compCache := mocks.NewMockcacheRepo(ctrl)

	uc := moderationstatus.NewUsecase(
		vacRepo,
		compRepo,
		outboxWriter,
		&fakeTxManager{},
		vacCache,
		pubVacCache,
		compCache,
	)

	req := &moderationstatus.Request{
		ID:     uuid.New(),
		Status: vacancy.ModerationStatusOK,
	}

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleAdmin,
	}

	compID := uuid.New()

	vacRepo.
		EXPECT().
		UpdateModerationStatus(
			gomock.Any(),
			req.ID,
			req.Status,
			gomock.Any(),
		).
		Return(&moderationstatus.UpdateModerationResult{
			CompanyID:           compID,
			VacancyStatus:       vacancy.StatusPublished,
			OldModerationStatus: vacancy.ModerationStatusHidden,
		}, nil)

	compRepo.
		EXPECT().
		IncrementOpenVacancies(gomock.Any(), compID).
		Return(nil)

	outboxWriter.
		EXPECT().
		WriteVacancyModerationUpdated(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, ev vacancy.ModerationUpdatedEvent) error {
			require.Equal(t, req.ID, ev.VacancyID)
			require.Equal(t, req.Status, ev.ModerationStatus)
			require.NotEqual(t, uuid.Nil, ev.EventID)

			return nil
		})

	vacCache.EXPECT().Del(gomock.Any(), req.ID)

	pubVacCache.EXPECT().Del(gomock.Any(), req.ID)

	compCache.EXPECT().Del(gomock.Any(), compID)

	err := uc.Execute(context.Background(), req, ident)

	require.NoError(t, err)
}

func TestUsecase_Execute_ArchivedVacancy_DoesNotUpdateCompanyCount(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	vacRepo := mocks.NewMockvacancyRepo(ctrl)
	compRepo := mocks.NewMockcompanyRepo(ctrl)
	outboxWriter := mocks.NewMockoutboxWriter(ctrl)

	vacCache := mocks.NewMockcacheRepo(ctrl)
	pubVacCache := mocks.NewMockcacheRepo(ctrl)
	compCache := mocks.NewMockcacheRepo(ctrl)

	uc := moderationstatus.NewUsecase(
		vacRepo,
		compRepo,
		outboxWriter,
		&fakeTxManager{},
		vacCache,
		pubVacCache,
		compCache,
	)

	req := &moderationstatus.Request{
		ID:     uuid.New(),
		Status: vacancy.ModerationStatusHidden,
	}

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleAdmin,
	}

	compID := uuid.New()

	vacRepo.
		EXPECT().
		UpdateModerationStatus(
			gomock.Any(),
			req.ID,
			req.Status,
			gomock.Any(),
		).
		Return(&moderationstatus.UpdateModerationResult{
			CompanyID:           compID,
			VacancyStatus:       vacancy.StatusArchived,
			OldModerationStatus: vacancy.ModerationStatusOK,
		}, nil)

	outboxWriter.
		EXPECT().
		WriteVacancyModerationUpdated(gomock.Any(), gomock.Any()).
		Return(nil)

	vacCache.
		EXPECT().
		Del(gomock.Any(), req.ID)

	pubVacCache.
		EXPECT().
		Del(gomock.Any(), req.ID)

	err := uc.Execute(context.Background(), req, ident)

	require.NoError(t, err)
}

func TestUsecase_Execute_StatusNotChanged(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	vacRepo := mocks.NewMockvacancyRepo(ctrl)

	uc := moderationstatus.NewUsecase(
		vacRepo,
		mocks.NewMockcompanyRepo(ctrl),
		mocks.NewMockoutboxWriter(ctrl),
		&fakeTxManager{},
		mocks.NewMockcacheRepo(ctrl),
		mocks.NewMockcacheRepo(ctrl),
		mocks.NewMockcacheRepo(ctrl),
	)

	req := &moderationstatus.Request{
		ID:     uuid.New(),
		Status: vacancy.ModerationStatusOK,
	}

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleAdmin,
	}

	vacRepo.
		EXPECT().
		UpdateModerationStatus(
			gomock.Any(),
			req.ID,
			req.Status,
			gomock.Any(),
		).
		Return(&moderationstatus.UpdateModerationResult{
			CompanyID:           uuid.New(),
			VacancyStatus:       vacancy.StatusPublished,
			OldModerationStatus: vacancy.ModerationStatusOK,
		}, nil)

	err := uc.Execute(context.Background(), req, ident)

	require.NoError(t, err)
}

func TestUsecase_Execute_InvalidStatus(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	uc := moderationstatus.NewUsecase(
		mocks.NewMockvacancyRepo(ctrl),
		mocks.NewMockcompanyRepo(ctrl),
		mocks.NewMockoutboxWriter(ctrl),
		&fakeTxManager{},
		mocks.NewMockcacheRepo(ctrl),
		mocks.NewMockcacheRepo(ctrl),
		mocks.NewMockcacheRepo(ctrl),
	)

	req := &moderationstatus.Request{
		ID:     uuid.New(),
		Status: vacancy.ModerationStatus("invalid"),
	}

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleAdmin,
	}

	err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, vacancy.ErrInvalidModerationStatus)
}

func TestUsecase_Execute_AdminRoleRequired(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	uc := moderationstatus.NewUsecase(
		mocks.NewMockvacancyRepo(ctrl),
		mocks.NewMockcompanyRepo(ctrl),
		mocks.NewMockoutboxWriter(ctrl),
		&fakeTxManager{},
		mocks.NewMockcacheRepo(ctrl),
		mocks.NewMockcacheRepo(ctrl),
		mocks.NewMockcacheRepo(ctrl),
	)

	req := &moderationstatus.Request{
		ID:     uuid.New(),
		Status: vacancy.ModerationStatusOK,
	}

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleHR,
	}

	err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, identity.ErrAdminRoleRequired)
}

func TestUsecase_Execute_UpdateModerationError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	vacRepo := mocks.NewMockvacancyRepo(ctrl)

	uc := moderationstatus.NewUsecase(
		vacRepo,
		mocks.NewMockcompanyRepo(ctrl),
		mocks.NewMockoutboxWriter(ctrl),
		&fakeTxManager{},
		mocks.NewMockcacheRepo(ctrl),
		mocks.NewMockcacheRepo(ctrl),
		mocks.NewMockcacheRepo(ctrl),
	)

	req := &moderationstatus.Request{
		ID:     uuid.New(),
		Status: vacancy.ModerationStatusOK,
	}

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleAdmin,
	}

	expectedErr := errors.New("update failed")

	vacRepo.
		EXPECT().
		UpdateModerationStatus(
			gomock.Any(),
			req.ID,
			req.Status,
			gomock.Any(),
		).
		Return(nil, expectedErr)

	err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, expectedErr)
}

func TestUsecase_Execute_DecrementOpenVacanciesError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	vacRepo := mocks.NewMockvacancyRepo(ctrl)
	compRepo := mocks.NewMockcompanyRepo(ctrl)

	uc := moderationstatus.NewUsecase(
		vacRepo,
		compRepo,
		mocks.NewMockoutboxWriter(ctrl),
		&fakeTxManager{},
		mocks.NewMockcacheRepo(ctrl),
		mocks.NewMockcacheRepo(ctrl),
		mocks.NewMockcacheRepo(ctrl),
	)

	req := &moderationstatus.Request{
		ID:     uuid.New(),
		Status: vacancy.ModerationStatusHidden,
	}

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleAdmin,
	}

	compID := uuid.New()

	expectedErr := errors.New("decrement failed")

	vacRepo.
		EXPECT().
		UpdateModerationStatus(
			gomock.Any(),
			req.ID,
			req.Status,
			gomock.Any(),
		).
		Return(&moderationstatus.UpdateModerationResult{
			CompanyID:           compID,
			VacancyStatus:       vacancy.StatusPublished,
			OldModerationStatus: vacancy.ModerationStatusOK,
		}, nil)

	compRepo.
		EXPECT().
		DecrementOpenVacancies(gomock.Any(), compID).
		Return(expectedErr)

	err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, expectedErr)
}
