package remove_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/remove"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/remove/mocks"
)

type fakeTxManager struct {
	called bool
}

func (f *fakeTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	f.called = true
	return fn(ctx)
}

type testDeps struct {
	compRepo  *mocks.MockCompanyRepo
	memRepo   *mocks.MockCompMemberRepo
	outbox    *mocks.MockoutboxWriter
	cache     *mocks.MockCacheRepo
	txManager *fakeTxManager
}

func setup(t *testing.T) *testDeps {
	ctrl := gomock.NewController(t)

	return &testDeps{
		compRepo:  mocks.NewMockCompanyRepo(ctrl),
		memRepo:   mocks.NewMockCompMemberRepo(ctrl),
		outbox:    mocks.NewMockoutboxWriter(ctrl),
		cache:     mocks.NewMockCacheRepo(ctrl),
		txManager: new(fakeTxManager),
	}
}

func NewUC(deps *testDeps) *remove.Usecase {
	return remove.NewUsecase(deps.compRepo, deps.memRepo, deps.outbox, deps.txManager, deps.cache)
}

type deletedEventMatcher struct {
	expected company.DeletedEvent
}

func (m deletedEventMatcher) Matches(x any) bool {
	ev, ok := x.(company.DeletedEvent)
	if !ok {
		return false
	}

	return ev.CompanyID == m.expected.CompanyID
}

func (m deletedEventMatcher) String() string {
	return "match company deleted event"
}

func TestUsecase_Execute_Success_HRAdmin(t *testing.T) {
	deps := setup(t)
	compID := uuid.New()

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleHR,
	}

	deps.memRepo.EXPECT().
		Get(gomock.Any(), ident.UserID, compID).
		Return(&member.CompanyMember{Role: member.CompanyRoleAdmin}, nil)

	deps.compRepo.EXPECT().Delete(gomock.Any(), compID).Return(nil)

	deps.outbox.EXPECT().WriteCompanyDeleted(gomock.Any(), deletedEventMatcher{
		expected: company.DeletedEvent{CompanyID: compID},
	})

	deps.cache.EXPECT().Del(gomock.Any(), compID)

	uc := NewUC(deps)

	err := uc.Execute(context.Background(), compID, ident)
	require.NoError(t, err)
}

func TestUsecase_Execute_Success_PlatformAdmin(t *testing.T) {
	deps := setup(t)
	compID := uuid.New()

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleAdmin}

	deps.compRepo.EXPECT().Delete(gomock.Any(), compID).Return(nil)

	deps.outbox.EXPECT().WriteCompanyDeleted(gomock.Any(), deletedEventMatcher{
		expected: company.DeletedEvent{CompanyID: compID},
	})

	deps.cache.EXPECT().Del(gomock.Any(), compID)

	uc := NewUC(deps)

	err := uc.Execute(context.Background(), compID, ident)
	require.NoError(t, err)
}

func TestUsecase_Execute_NotFound(t *testing.T) {
	deps := setup(t)
	compID := uuid.New()

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleAdmin,
	}

	deps.compRepo.EXPECT().
		Delete(gomock.Any(), compID).
		Return(company.ErrCompanyNotFound)

	uc := NewUC(deps)

	err := uc.Execute(context.Background(), compID, ident)

	require.ErrorIs(t, err, company.ErrCompanyNotFound)
}

func TestUsecase_Execute_AuthError_Role(t *testing.T) {
	deps := setup(t)
	compID := uuid.New()

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleCandidate}

	uc := NewUC(deps)

	err := uc.Execute(context.Background(), compID, ident)

	require.ErrorIs(t, err, identity.ErrInsufficientRole)
}

func TestUsecase_Execute_AuthError_NotMember(t *testing.T) {
	deps := setup(t)
	compID := uuid.New()

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleHR,
	}

	deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, compID).
		Return(nil, member.ErrCompanyMemberNotFound)

	uc := NewUC(deps)

	err := uc.Execute(context.Background(), compID, ident)

	require.ErrorIs(t, err, member.ErrCompanyMemberRequired)
}
