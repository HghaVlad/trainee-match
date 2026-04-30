package publish_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/publish"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/publish/mocks"
)

type fakeTxManager struct {
	called bool
}

func (f *fakeTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	f.called = true
	return fn(ctx)
}

type testDeps struct {
	memRepo      *mocks.MockCompMemberRepo
	vacRepo      *mocks.MockVacancyRepo
	compRepo     *mocks.MockCompanyRepo
	outboxWriter *mocks.MockoutboxWriter
	vacCache     *mocks.MockCacheRepo
	compCache    *mocks.MockCacheRepo
	txManager    *fakeTxManager
}

func setup(t *testing.T) *testDeps {
	ctrl := gomock.NewController(t)

	return &testDeps{
		memRepo:      mocks.NewMockCompMemberRepo(ctrl),
		vacRepo:      mocks.NewMockVacancyRepo(ctrl),
		compRepo:     mocks.NewMockCompanyRepo(ctrl),
		outboxWriter: mocks.NewMockoutboxWriter(ctrl),
		vacCache:     mocks.NewMockCacheRepo(ctrl),
		compCache:    mocks.NewMockCacheRepo(ctrl),
		txManager:    new(fakeTxManager),
	}
}

func NewUC(deps *testDeps) *publish.Usecase {
	return publish.NewUsecase(
		deps.vacRepo,
		deps.compRepo,
		deps.memRepo,
		deps.outboxWriter,
		deps.txManager,
		deps.vacCache,
		deps.compCache,
	)
}

type pubEventMatcher struct {
	expectedEv vacancy.PublishedEvent
}

func (m pubEventMatcher) Matches(x any) bool {
	ev, ok := x.(vacancy.PublishedEvent)
	if !ok {
		return false
	}

	return ev.VacancyID == m.expectedEv.VacancyID && ev.CompanyID == m.expectedEv.CompanyID &&
		ev.Title == m.expectedEv.Title && ev.CompanyName == m.expectedEv.CompanyName
}

func (m pubEventMatcher) String() string {
	return "matches vacancy.PublishedEvent"
}

func TestUsecase_Execute_PublishesDraftVacancy(t *testing.T) {
	deps := setup(t)

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}
	compID := uuid.New()
	vacID := uuid.New()

	eventView := publish.PublishedEventView{
		ID:          vacID,
		Title:       "title",
		CompanyID:   compID,
		CompanyName: "name",
		Status:      vacancy.StatusDraft,
	}

	event := vacancy.PublishedEvent{
		VacancyID:   vacID,
		CompanyID:   compID,
		Title:       "title",
		CompanyName: "name",
	}

	deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, compID).
		Return(&member.CompanyMember{UserID: ident.UserID, CompanyID: compID, Role: member.CompanyRoleRecruiter}, nil)

	deps.vacRepo.EXPECT().GetPublishedEventView(gomock.Any(), vacID, compID).Return(&eventView, nil)

	deps.compRepo.EXPECT().IncrementOpenVacancies(gomock.Any(), compID).Return(nil)

	deps.vacRepo.EXPECT().Publish(gomock.Any(), vacID, compID).Return(nil)

	deps.outboxWriter.EXPECT().WriteVacancyPublished(gomock.Any(), pubEventMatcher{expectedEv: event}).Return(nil)

	deps.vacCache.EXPECT().Del(gomock.Any(), vacID)
	deps.compCache.EXPECT().Del(gomock.Any(), compID)

	uc := NewUC(deps)

	err := uc.Execute(context.Background(), compID, vacID, ident)

	require.NoError(t, err)
	assert.True(t, deps.txManager.called)
}

func TestUsecase_Execute_Alreadypublish_NoOp(t *testing.T) {
	deps := setup(t)

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}
	compID := uuid.New()
	vacID := uuid.New()

	deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, compID).
		Return(&member.CompanyMember{UserID: ident.UserID, CompanyID: compID, Role: member.CompanyRoleRecruiter}, nil)

	deps.vacRepo.EXPECT().GetPublishedEventView(gomock.Any(), vacID, compID).
		Return(&publish.PublishedEventView{ID: vacID, CompanyID: compID, Status: vacancy.StatusPublished}, nil)

	uc := NewUC(deps)

	err := uc.Execute(t.Context(), compID, vacID, ident)

	require.NoError(t, err)
	require.True(t, deps.txManager.called)
}

func TestUsecase_Execute_AuthErr(t *testing.T) {
	compID := uuid.New()
	vacID := uuid.New()

	t.Run("hr role required", func(t *testing.T) {
		deps := setup(t)

		uc := NewUC(deps)

		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleCandidate}
		err := uc.Execute(context.Background(), compID, vacID, ident)

		require.ErrorIs(t, err, identity.ErrHrRoleRequired)
	})

	t.Run("company member required", func(t *testing.T) {
		deps := setup(t)

		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

		deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, compID).
			Return(nil, member.ErrCompanyMemberNotFound)

		uc := NewUC(deps)

		err := uc.Execute(t.Context(), compID, vacID, ident)

		require.ErrorIs(t, err, member.ErrCompanyMemberRequired)
	})
}

func TestUsecase_Execute_RepoErr(t *testing.T) {
	t.Run("get vacancy", func(t *testing.T) {
		deps := setup(t)

		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}
		compID := uuid.New()
		vacID := uuid.New()

		deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, compID).
			Return(&member.CompanyMember{UserID: ident.UserID, CompanyID: compID}, nil)

		deps.vacRepo.EXPECT().GetPublishedEventView(gomock.Any(), vacID, compID).
			Return(nil, vacancy.ErrVacancyNotFound)

		uc := NewUC(deps)

		err := uc.Execute(context.Background(), compID, vacID, ident)

		require.ErrorIs(t, err, vacancy.ErrVacancyNotFound)
	})

	t.Run("publish", func(t *testing.T) {
		deps := setup(t)

		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}
		compID := uuid.New()
		vacID := uuid.New()

		deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, compID).
			Return(&member.CompanyMember{}, nil)

		deps.vacRepo.EXPECT().GetPublishedEventView(gomock.Any(), vacID, compID).
			Return(&publish.PublishedEventView{ID: vacID, CompanyID: compID, Status: vacancy.StatusDraft}, nil)

		deps.vacRepo.EXPECT().Publish(gomock.Any(), vacID, compID).
			Return(errors.New("db err"))

		uc := NewUC(deps)

		err := uc.Execute(context.Background(), compID, vacID, ident)

		require.Error(t, err)
	})
}
