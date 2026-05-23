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
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/views"
)

type fakeTxManager struct{}

func (f *fakeTxManager) WithinTx(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	return fn(ctx)
}

type testDeps struct {
	vacRepo   *mocks.MockvacancyRepo
	compRepo  *mocks.MockcompanyRepo
	outbox    *mocks.MockoutboxWriter
	search    *mocks.MocksearchRepo
	vacCache  *mocks.MockcacheRepo
	pubCache  *mocks.MockcacheRepo
	compCache *mocks.MockcacheRepo
	txManager *fakeTxManager
}

func setup(t *testing.T) *testDeps {
	ctrl := gomock.NewController(t)

	return &testDeps{
		vacRepo:   mocks.NewMockvacancyRepo(ctrl),
		compRepo:  mocks.NewMockcompanyRepo(ctrl),
		outbox:    mocks.NewMockoutboxWriter(ctrl),
		search:    mocks.NewMocksearchRepo(ctrl),
		vacCache:  mocks.NewMockcacheRepo(ctrl),
		pubCache:  mocks.NewMockcacheRepo(ctrl),
		compCache: mocks.NewMockcacheRepo(ctrl),
		txManager: new(fakeTxManager),
	}
}

func newUC(deps *testDeps) *moderationstatus.Usecase {
	return moderationstatus.NewUsecase(
		deps.vacRepo,
		deps.compRepo,
		deps.outbox,
		deps.txManager,
		deps.search,
		deps.vacCache,
		deps.pubCache,
		deps.compCache,
	)
}

func adminIdentity() *identity.Identity {
	return &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleAdmin,
	}
}

func TestUsecase_Execute_Success_StatusChanged_Hidden(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUC(deps)

	vacID := uuid.New()
	compID := uuid.New()

	req := &moderationstatus.Request{
		ID:     vacID,
		Status: vacancy.ModerationStatusHidden,
	}

	searchView := &views.VacancySearch{
		ID:               vacID,
		CompanyID:        compID,
		ModerationStatus: vacancy.ModerationStatusHidden,
		Status:           vacancy.StatusPublished,
	}

	deps.vacRepo.EXPECT().
		UpdateModerationStatus(
			gomock.Any(),
			vacID,
			vacancy.ModerationStatusHidden,
			gomock.Any(),
		).
		Return(&moderationstatus.UpdateModerationResult{
			CompanyID:           compID,
			OldModerationStatus: vacancy.ModerationStatusOK,
			VacancyStatus:       vacancy.StatusPublished,
		}, nil)

	deps.compRepo.EXPECT().
		DecrementOpenVacancies(gomock.Any(), compID).
		Return(nil)

	deps.outbox.EXPECT().
		WriteVacancyModerationUpdated(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, ev vacancy.ModerationUpdatedEvent) error {
			require.Equal(t, vacID, ev.VacancyID)
			require.Equal(t, vacancy.ModerationStatusHidden, ev.ModerationStatus)
			require.NotEqual(t, uuid.Nil, ev.EventID)

			return nil
		})

	deps.vacCache.EXPECT().
		Del(gomock.Any(), vacID)

	deps.pubCache.EXPECT().
		Del(gomock.Any(), vacID)

	deps.vacRepo.EXPECT().
		GetSearchView(gomock.Any(), vacID).
		Return(searchView, nil)

	deps.search.EXPECT().
		Index(gomock.Any(), *searchView).
		Return(nil)

	deps.compCache.EXPECT().
		Del(gomock.Any(), compID)

	err := uc.Execute(context.Background(), req, adminIdentity())
	require.NoError(t, err)
}

func TestUsecase_Execute_Success_StatusChanged_OK(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUC(deps)

	vacID := uuid.New()
	compID := uuid.New()

	req := &moderationstatus.Request{
		ID:     vacID,
		Status: vacancy.ModerationStatusOK,
	}

	searchView := &views.VacancySearch{
		ID:               vacID,
		CompanyID:        compID,
		ModerationStatus: vacancy.ModerationStatusOK,
		Status:           vacancy.StatusPublished,
	}

	deps.vacRepo.EXPECT().
		UpdateModerationStatus(
			gomock.Any(),
			vacID,
			vacancy.ModerationStatusOK,
			gomock.Any(),
		).
		Return(&moderationstatus.UpdateModerationResult{
			CompanyID:           compID,
			OldModerationStatus: vacancy.ModerationStatusHidden,
			VacancyStatus:       vacancy.StatusPublished,
		}, nil)

	deps.compRepo.EXPECT().
		IncrementOpenVacancies(gomock.Any(), compID).
		Return(nil)

	deps.outbox.EXPECT().
		WriteVacancyModerationUpdated(gomock.Any(), gomock.Any()).
		Return(nil)

	deps.vacCache.EXPECT().
		Del(gomock.Any(), vacID)

	deps.pubCache.EXPECT().
		Del(gomock.Any(), vacID)

	deps.vacRepo.EXPECT().
		GetSearchView(gomock.Any(), vacID).
		Return(searchView, nil)

	deps.search.EXPECT().
		Index(gomock.Any(), *searchView).
		Return(nil)

	deps.compCache.EXPECT().
		Del(gomock.Any(), compID)

	err := uc.Execute(context.Background(), req, adminIdentity())
	require.NoError(t, err)
}

func TestUsecase_Execute_Success_StatusNotChanged(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUC(deps)

	vacID := uuid.New()
	compID := uuid.New()

	req := &moderationstatus.Request{
		ID:     vacID,
		Status: vacancy.ModerationStatusOK,
	}

	deps.vacRepo.EXPECT().
		UpdateModerationStatus(
			gomock.Any(),
			vacID,
			vacancy.ModerationStatusOK,
			gomock.Any(),
		).
		Return(&moderationstatus.UpdateModerationResult{
			CompanyID:           compID,
			OldModerationStatus: vacancy.ModerationStatusOK,
			VacancyStatus:       vacancy.StatusPublished,
		}, nil)

	err := uc.Execute(context.Background(), req, adminIdentity())
	require.NoError(t, err)
}

func TestUsecase_Execute_Success_ArchivedVacancy_NoCompanyUpdate(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUC(deps)

	vacID := uuid.New()
	compID := uuid.New()

	req := &moderationstatus.Request{
		ID:     vacID,
		Status: vacancy.ModerationStatusHidden,
	}

	searchView := &views.VacancySearch{
		ID:               vacID,
		CompanyID:        compID,
		ModerationStatus: vacancy.ModerationStatusHidden,
		Status:           vacancy.StatusArchived,
	}

	deps.vacRepo.EXPECT().
		UpdateModerationStatus(
			gomock.Any(),
			vacID,
			vacancy.ModerationStatusHidden,
			gomock.Any(),
		).
		Return(&moderationstatus.UpdateModerationResult{
			CompanyID:           compID,
			OldModerationStatus: vacancy.ModerationStatusOK,
			VacancyStatus:       vacancy.StatusArchived,
		}, nil)

	deps.outbox.EXPECT().
		WriteVacancyModerationUpdated(gomock.Any(), gomock.Any()).
		Return(nil)

	deps.vacCache.EXPECT().
		Del(gomock.Any(), vacID)

	deps.pubCache.EXPECT().
		Del(gomock.Any(), vacID)

	deps.vacRepo.EXPECT().
		GetSearchView(gomock.Any(), vacID).
		Return(searchView, nil)

	deps.search.EXPECT().
		Index(gomock.Any(), *searchView).
		Return(nil)

	err := uc.Execute(context.Background(), req, adminIdentity())
	require.NoError(t, err)
}

func TestUsecase_Execute_AuthError(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUC(deps)

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

	deps := setup(t)
	uc := newUC(deps)

	vacID := uuid.New()

	req := &moderationstatus.Request{
		ID:     vacID,
		Status: vacancy.ModerationStatusHidden,
	}

	expectedErr := errors.New("update moderation error")

	deps.vacRepo.EXPECT().
		UpdateModerationStatus(
			gomock.Any(),
			vacID,
			vacancy.ModerationStatusHidden,
			gomock.Any(),
		).
		Return(nil, expectedErr)

	err := uc.Execute(context.Background(), req, adminIdentity())

	require.ErrorIs(t, err, expectedErr)
}

func TestUsecase_Execute_CompanyUpdateError(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUC(deps)

	vacID := uuid.New()
	compID := uuid.New()

	req := &moderationstatus.Request{
		ID:     vacID,
		Status: vacancy.ModerationStatusHidden,
	}

	expectedErr := errors.New("company update error")

	deps.vacRepo.EXPECT().
		UpdateModerationStatus(
			gomock.Any(),
			vacID,
			vacancy.ModerationStatusHidden,
			gomock.Any(),
		).
		Return(&moderationstatus.UpdateModerationResult{
			CompanyID:           compID,
			OldModerationStatus: vacancy.ModerationStatusOK,
			VacancyStatus:       vacancy.StatusPublished,
		}, nil)

	deps.compRepo.EXPECT().
		DecrementOpenVacancies(gomock.Any(), compID).
		Return(expectedErr)

	err := uc.Execute(context.Background(), req, adminIdentity())

	require.ErrorIs(t, err, expectedErr)
}

func TestUsecase_Execute_OutboxError(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUC(deps)

	vacID := uuid.New()
	compID := uuid.New()

	req := &moderationstatus.Request{
		ID:     vacID,
		Status: vacancy.ModerationStatusHidden,
	}

	expectedErr := errors.New("outbox error")

	deps.vacRepo.EXPECT().
		UpdateModerationStatus(
			gomock.Any(),
			vacID,
			vacancy.ModerationStatusHidden,
			gomock.Any(),
		).
		Return(&moderationstatus.UpdateModerationResult{
			CompanyID:           compID,
			OldModerationStatus: vacancy.ModerationStatusOK,
			VacancyStatus:       vacancy.StatusPublished,
		}, nil)

	deps.compRepo.EXPECT().
		DecrementOpenVacancies(gomock.Any(), compID).
		Return(nil)

	deps.outbox.EXPECT().
		WriteVacancyModerationUpdated(gomock.Any(), gomock.Any()).
		Return(expectedErr)

	err := uc.Execute(context.Background(), req, adminIdentity())

	require.ErrorIs(t, err, expectedErr)
}
