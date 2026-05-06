package add_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/add"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/add/mocks"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/projection/userhr"
)

type fakeTxManager struct {
	called bool
}

func (f *fakeTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	f.called = true
	return fn(ctx)
}

type testDeps struct {
	memRepo    *mocks.MockcompanyMemberRepo
	hrProjRepo *mocks.MockhrProjRepo
	outbox     *mocks.MockoutboxWriter
	txManager  *fakeTxManager
}

func setup(t *testing.T) *testDeps {
	ctrl := gomock.NewController(t)
	return &testDeps{
		memRepo:    mocks.NewMockcompanyMemberRepo(ctrl),
		hrProjRepo: mocks.NewMockhrProjRepo(ctrl),
		outbox:     mocks.NewMockoutboxWriter(ctrl),
		txManager:  new(fakeTxManager),
	}
}

type memberMatcher struct {
	expected *member.CompanyMember
}

func (m memberMatcher) Matches(x any) bool {
	mem, ok := x.(*member.CompanyMember)
	if !ok {
		return false
	}

	return m.expected.UserID == mem.UserID && m.expected.CompanyID == mem.CompanyID &&
		m.expected.Role == mem.Role
}

func (m memberMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.expected)
}

type memberAddedEvMatcher struct {
	expected member.AddedEvent
}

func (m memberAddedEvMatcher) Matches(x any) bool {
	ev, ok := x.(member.AddedEvent)
	if !ok {
		return false
	}

	return m.expected.UserID == ev.UserID && m.expected.CompanyID == ev.CompanyID &&
		m.expected.Role == ev.Role
}

func (m memberAddedEvMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.expected)
}

func NewUC(deps *testDeps) *add.Usecase {
	return add.NewUsecase(deps.memRepo, deps.hrProjRepo, deps.outbox, deps.txManager)
}

func TestUsecase_ExecuteOK(t *testing.T) {
	deps := setup(t)

	compID := uuid.New()
	usID := uuid.New()
	usname := "JohnPork360"
	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}
	req := &add.Request{
		Username:  usname,
		CompanyID: compID,
		Role:      member.CompanyRoleRecruiter,
	}

	mem := &member.CompanyMember{
		UserID:    usID,
		CompanyID: req.CompanyID,
		Role:      req.Role,
	}

	hrProj := &userhr.Projection{
		UserID:    usID,
		Username:  usname,
		Email:     "usname@mail.com",
		CreatedAt: time.Now().UTC(),
	}

	memEv := member.AddedEvent{UserID: usID, CompanyID: req.CompanyID, Role: req.Role}

	deps.hrProjRepo.EXPECT().GetByUsername(gomock.Any(), usname).
		Return(hrProj, nil)

	deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, compID).
		Return(&member.CompanyMember{Role: member.CompanyRoleAdmin}, nil)

	deps.memRepo.EXPECT().Create(gomock.Any(), memberMatcher{expected: mem}).
		Return(nil)

	deps.outbox.EXPECT().WriteCompanyMemberAdded(gomock.Any(), memberAddedEvMatcher{memEv}).Return(nil)

	uc := NewUC(deps)

	err := uc.Execute(t.Context(), req, ident)

	require.NoError(t, err)
}

func TestUsecase_ExecuteUsProjNotFound(t *testing.T) {
	usname := "JohnPork360"
	req := &add.Request{
		CompanyID: uuid.New(),
		Username:  usname,
		Role:      member.CompanyRoleRecruiter,
	}
	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	deps := setup(t)

	deps.hrProjRepo.EXPECT().GetByUsername(gomock.Any(), usname).
		Return(nil, userhr.ErrNotFound)

	uc := NewUC(deps)

	err := uc.Execute(t.Context(), req, ident)

	require.ErrorIs(t, err, userhr.ErrNotFound)
}

func TestUsecase_ExecuteAuthErr(t *testing.T) {
	usname := "JohnPork360"
	req := &add.Request{
		CompanyID: uuid.New(),
		Username:  usname,
		Role:      member.CompanyRoleRecruiter,
	}

	t.Run("global hr role required", func(t *testing.T) {
		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleCandidate}

		deps := setup(t)

		uc := NewUC(deps)

		err := uc.Execute(t.Context(), req, ident)

		require.ErrorIs(t, err, identity.ErrHrRoleRequired)
	})

	t.Run("company member required", func(t *testing.T) {
		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

		deps := setup(t)

		deps.hrProjRepo.EXPECT().GetByUsername(gomock.Any(), usname).
			Return(&userhr.Projection{UserID: uuid.New(), Username: usname}, nil)

		deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, req.CompanyID).
			Return(nil, member.ErrCompanyMemberNotFound)

		uc := NewUC(deps)

		err := uc.Execute(t.Context(), req, ident)

		require.ErrorIs(t, err, member.ErrCompanyMemberRequired)
	})

	t.Run("admin company role required", func(t *testing.T) {
		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

		deps := setup(t)

		deps.hrProjRepo.EXPECT().GetByUsername(gomock.Any(), usname).
			Return(&userhr.Projection{UserID: uuid.New(), Username: usname}, nil)

		deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, req.CompanyID).
			Return(&member.CompanyMember{Role: member.CompanyRoleRecruiter}, nil)

		uc := NewUC(deps)

		err := uc.Execute(t.Context(), req, ident)

		require.ErrorIs(t, err, member.ErrInsufficientRoleInCompany)
	})
}

func TestUsecase_ExecuteValidationAndRepoErr(t *testing.T) {
	usname := "JohnPork360"

	t.Run("invalid role", func(t *testing.T) {
		deps := setup(t)
		uc := NewUC(deps)

		req := &add.Request{
			CompanyID: uuid.New(),
			Username:  usname,
			Role:      "invalid",
		}
		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

		err := uc.Execute(t.Context(), req, ident)

		require.ErrorIs(t, err, member.ErrInvalidCompanyMemberRole)
	})

	t.Run("repo err", func(t *testing.T) {
		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}
		req := &add.Request{
			CompanyID: uuid.New(),
			Username:  usname,
			Role:      member.CompanyRoleAdmin,
		}

		deps := setup(t)

		deps.hrProjRepo.EXPECT().GetByUsername(gomock.Any(), usname).
			Return(&userhr.Projection{UserID: uuid.New(), Username: usname}, nil)

		deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, req.CompanyID).
			Return(&member.CompanyMember{Role: member.CompanyRoleAdmin}, nil)

		deps.memRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
			Return(member.ErrCompanyMemberAlreadyExists)

		uc := NewUC(deps)

		err := uc.Execute(t.Context(), req, ident)

		require.ErrorIs(t, err, member.ErrCompanyMemberAlreadyExists)
	})
}
