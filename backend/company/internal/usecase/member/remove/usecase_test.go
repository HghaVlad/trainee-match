package remove_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/remove"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/remove/mocks"
)

type fakeTxManager struct {
	called bool
}

func (f *fakeTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	f.called = true
	return fn(ctx)
}

type testDeps struct {
	memRepo   *mocks.MockCompanyMemberRepo
	outbox    *mocks.MockoutboxWriter
	txManager *fakeTxManager
}

func setup(t *testing.T) *testDeps {
	ctrl := gomock.NewController(t)
	return &testDeps{
		memRepo:   mocks.NewMockCompanyMemberRepo(ctrl),
		outbox:    mocks.NewMockoutboxWriter(ctrl),
		txManager: new(fakeTxManager),
	}
}

type memberRemovedEvMatcher struct {
	expected member.RemovedEvent
}

func (m memberRemovedEvMatcher) Matches(x any) bool {
	ev, ok := x.(member.RemovedEvent)
	if !ok {
		return false
	}

	return m.expected.UserID == ev.UserID && m.expected.CompanyID == ev.CompanyID
}

func (m memberRemovedEvMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.expected)
}

func NewUC(deps *testDeps) *remove.Usecase {
	return remove.NewUsecase(deps.memRepo, deps.outbox, deps.txManager)
}

func TestUsecase_ExecuteOK(t *testing.T) {
	deps := setup(t)

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}
	companyID := uuid.New()
	userID := uuid.New()

	expectedEv := member.RemovedEvent{UserID: userID, CompanyID: companyID}

	deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, companyID).
		Return(&member.CompanyMember{Role: member.CompanyRoleAdmin}, nil)

	deps.memRepo.EXPECT().Delete(gomock.Any(), userID, companyID).Return(nil)

	deps.outbox.EXPECT().
		WriteCompanyMemberRemoved(gomock.Any(), memberRemovedEvMatcher{expected: expectedEv}).Return(nil)

	uc := NewUC(deps)

	err := uc.Execute(t.Context(), companyID, userID, ident)

	require.NoError(t, err)
}

func TestUsecase_ExecuteRemovesThemselvesWhenMoreThanOneAdmin(t *testing.T) {
	deps := setup(t)

	userID := uuid.New()
	ident := &identity.Identity{UserID: userID, Role: identity.RoleHR}
	companyID := uuid.New()

	expectedEv := member.RemovedEvent{UserID: userID, CompanyID: companyID}

	deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, companyID).
		Return(&member.CompanyMember{Role: member.CompanyRoleAdmin}, nil)

	deps.memRepo.EXPECT().
		GetCompanyRoleCount(gomock.Any(), companyID, member.CompanyRoleAdmin).
		Return(2, nil)

	deps.memRepo.EXPECT().Delete(gomock.Any(), userID, companyID).Return(nil)

	deps.outbox.EXPECT().
		WriteCompanyMemberRemoved(gomock.Any(), memberRemovedEvMatcher{expected: expectedEv}).Return(nil)

	uc := NewUC(deps)

	err := uc.Execute(t.Context(), companyID, userID, ident)

	require.NoError(t, err)
}

func TestUsecase_ExecuteCantRemovesThemselvesWhenTheOnlyAdmin(t *testing.T) {
	deps := setup(t)

	userID := uuid.New()
	ident := &identity.Identity{UserID: userID, Role: identity.RoleHR}
	companyID := uuid.New()

	deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, companyID).
		Return(&member.CompanyMember{Role: member.CompanyRoleAdmin}, nil)

	deps.memRepo.EXPECT().
		GetCompanyRoleCount(gomock.Any(), companyID, member.CompanyRoleAdmin).
		Return(1, nil)

	uc := NewUC(deps)

	err := uc.Execute(t.Context(), companyID, userID, ident)

	require.ErrorIs(t, err, member.ErrCantRemoveYourself)
}

func TestUsecase_ExecuteAuthErr(t *testing.T) {
	companyID := uuid.New()
	userID := uuid.New()

	t.Run("global hr role required", func(t *testing.T) {
		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleCandidate}

		deps := setup(t)

		uc := NewUC(deps)

		err := uc.Execute(t.Context(), companyID, userID, ident)

		require.ErrorIs(t, err, identity.ErrHrRoleRequired)
	})

	t.Run("company member required", func(t *testing.T) {
		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

		deps := setup(t)

		deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, companyID).
			Return(nil, member.ErrCompanyMemberNotFound)

		uc := NewUC(deps)

		err := uc.Execute(t.Context(), companyID, userID, ident)

		require.ErrorIs(t, err, member.ErrCompanyMemberRequired)
	})

	t.Run("admin company role required", func(t *testing.T) {
		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

		deps := setup(t)

		deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, companyID).
			Return(&member.CompanyMember{Role: member.CompanyRoleRecruiter}, nil)

		uc := NewUC(deps)

		err := uc.Execute(t.Context(), companyID, userID, ident)

		require.ErrorIs(t, err, member.ErrInsufficientRoleInCompany)
	})
}

func TestUsecase_ExecuteRepoErr(t *testing.T) {
	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}
	companyID := uuid.New()
	userID := uuid.New()

	deps := setup(t)

	deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, companyID).
		Return(&member.CompanyMember{Role: member.CompanyRoleAdmin}, nil)
	deps.memRepo.EXPECT().Delete(gomock.Any(), userID, companyID).
		Return(errors.New("db err"))

	uc := NewUC(deps)

	err := uc.Execute(context.Background(), companyID, userID, ident)

	require.Error(t, err)
}
