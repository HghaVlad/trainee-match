package archive_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/archive"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/archive/mocks"
)

type fakeTxManager struct {
	called bool
}

func (f *fakeTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	f.called = true
	return fn(ctx)
}

type archivedEventMatcher struct {
	expectedVacID uuid.UUID
}

func (m archivedEventMatcher) Matches(x any) bool {
	ev, ok := x.(vacancy.ArchivedEvent)
	if !ok {
		return false
	}
	return ev.VacancyID == m.expectedVacID
}

func (m archivedEventMatcher) String() string {
	return "matches ArchivedEvent with specific VacancyID"
}

type testDeps struct {
	vacRepo      *mocks.MockVacancyRepo
	compRepo     *mocks.MockCompanyRepo
	memRepo      *mocks.MockCompMemberRepo
	outboxWriter *mocks.MockoutboxWriter
	vacCache     *mocks.MockCacheRepo
	pubVacCache  *mocks.MockCacheRepo
	compCache    *mocks.MockCacheRepo
	txManager    *fakeTxManager
}

func setup(t *testing.T) *testDeps {
	ctrl := gomock.NewController(t)

	return &testDeps{
		vacRepo:      mocks.NewMockVacancyRepo(ctrl),
		compRepo:     mocks.NewMockCompanyRepo(ctrl),
		memRepo:      mocks.NewMockCompMemberRepo(ctrl),
		outboxWriter: mocks.NewMockoutboxWriter(ctrl),
		vacCache:     mocks.NewMockCacheRepo(ctrl),
		compCache:    mocks.NewMockCacheRepo(ctrl),
		pubVacCache:  mocks.NewMockCacheRepo(ctrl),
		txManager:    new(fakeTxManager),
	}
}

func NewUC(deps *testDeps) *archive.Usecase {
	return archive.NewUsecase(
		deps.vacRepo,
		deps.compRepo,
		deps.memRepo,
		deps.outboxWriter,
		deps.txManager,
		deps.vacCache,
		deps.pubVacCache,
		deps.compCache,
	)
}

func TestUsecase_Execute_ArchivesPublishedVacancy(t *testing.T) {
	ctx := t.Context()
	deps := setup(t)

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}
	compID := uuid.New()
	vacID := uuid.New()

	deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, compID).
		Return(&member.CompanyMember{UserID: ident.UserID}, nil)

	deps.vacRepo.EXPECT().ArchiveAndGetOldStatus(gomock.Any(), vacID, compID).
		Return(vacancy.StatusPublished, nil)

	deps.compRepo.EXPECT().DecrementOpenVacancies(gomock.Any(), compID).Return(nil)

	deps.outboxWriter.EXPECT().
		WriteVacancyArchived(gomock.Any(), archivedEventMatcher{vacID}).Return(nil)

	deps.vacCache.EXPECT().Del(gomock.Any(), vacID).Return()
	deps.pubVacCache.EXPECT().Del(gomock.Any(), vacID).Return()
	deps.compCache.EXPECT().Del(gomock.Any(), compID).Return()

	uc := NewUC(deps)

	err := uc.Execute(ctx, compID, vacID, ident)

	require.NoError(t, err)
	require.True(t, deps.txManager.called)
}

func TestUsecase_Execute_ArchivesDraftWithoutCounterUpdate(t *testing.T) {
	deps := setup(t)
	ctx := t.Context()

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}
	compID := uuid.New()
	vacID := uuid.New()

	deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, compID).
		Return(&member.CompanyMember{UserID: ident.UserID}, nil)

	deps.vacRepo.EXPECT().ArchiveAndGetOldStatus(gomock.Any(), vacID, compID).
		Return(vacancy.StatusDraft, nil)

	deps.vacCache.EXPECT().Del(gomock.Any(), vacID).Return()
	deps.pubVacCache.EXPECT().Del(gomock.Any(), vacID).Return()

	uc := NewUC(deps)

	err := uc.Execute(ctx, compID, vacID, ident)

	require.NoError(t, err)
	require.True(t, deps.txManager.called)
}

func TestUsecase_Execute_AlreadyArchived_NoOp(t *testing.T) {
	deps := setup(t)
	ctx := t.Context()

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}
	compID := uuid.New()
	vacID := uuid.New()

	deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, compID).
		Return(&member.CompanyMember{UserID: ident.UserID}, nil)

	deps.vacRepo.EXPECT().ArchiveAndGetOldStatus(gomock.Any(), vacID, compID).
		Return(vacancy.StatusArchived, nil)

	uc := NewUC(deps)

	err := uc.Execute(ctx, compID, vacID, ident)

	require.NoError(t, err)
	require.True(t, deps.txManager.called)
}

func TestUsecase_Execute_AuthErr(t *testing.T) {
	compID := uuid.New()
	vacID := uuid.New()

	t.Run("hr role required", func(t *testing.T) {
		deps := setup(t)
		ctx := t.Context()

		uc := NewUC(deps)

		iden := &identity.Identity{UserID: uuid.New(), Role: identity.RoleCandidate}
		err := uc.Execute(ctx, compID, vacID, iden)

		require.ErrorIs(t, err, identity.ErrHrRoleRequired)
	})

	t.Run("company member required", func(t *testing.T) {
		deps := setup(t)
		ctx := t.Context()

		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

		deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, compID).
			Return(nil, member.ErrCompanyMemberNotFound)

		uc := NewUC(deps)

		err := uc.Execute(ctx, compID, vacID, ident)

		require.ErrorIs(t, err, member.ErrCompanyMemberRequired)
	})
}

func TestUsecase_Execute_RepoErr(t *testing.T) {
	compID := uuid.New()
	vacID := uuid.New()
	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	t.Run("not found", func(t *testing.T) {
		deps := setup(t)
		ctx := t.Context()

		deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, compID).
			Return(&member.CompanyMember{UserID: ident.UserID}, nil)

		deps.vacRepo.EXPECT().ArchiveAndGetOldStatus(gomock.Any(), vacID, compID).
			Return(vacancy.Status(""), vacancy.ErrVacancyNotFound)

		uc := NewUC(deps)

		err := uc.Execute(ctx, compID, vacID, ident)

		require.ErrorIs(t, err, vacancy.ErrVacancyNotFound)
	})

	t.Run("decrement company vacancies company not found", func(t *testing.T) {
		deps := setup(t)
		ctx := t.Context()

		deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, compID).
			Return(&member.CompanyMember{UserID: ident.UserID}, nil)

		deps.vacRepo.EXPECT().ArchiveAndGetOldStatus(gomock.Any(), vacID, compID).
			Return(vacancy.StatusPublished, nil)

		deps.compRepo.EXPECT().DecrementOpenVacancies(gomock.Any(), compID).
			Return(company.ErrCompanyNotFound)

		uc := NewUC(deps)

		err := uc.Execute(ctx, compID, vacID, ident)

		require.ErrorIs(t, err, company.ErrCompanyNotFound)
	})
}
