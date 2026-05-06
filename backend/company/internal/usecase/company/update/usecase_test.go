package update_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/update"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/update/mocks"
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
	compCache *mocks.MockCacheRepo
	txManager *fakeTxManager
}

func setup(t *testing.T) *testDeps {
	ctrl := gomock.NewController(t)

	return &testDeps{
		compRepo:  mocks.NewMockCompanyRepo(ctrl),
		memRepo:   mocks.NewMockCompMemberRepo(ctrl),
		outbox:    mocks.NewMockoutboxWriter(ctrl),
		compCache: mocks.NewMockCacheRepo(ctrl),
		txManager: new(fakeTxManager),
	}
}

func NewUC(deps *testDeps) *update.Usecase {
	return update.NewUsecase(deps.compRepo, deps.memRepo, deps.outbox, deps.txManager, deps.compCache)
}

type eventMatcher struct {
	expectedEv company.UpdatedEvent
}

func (m eventMatcher) Matches(x any) bool {
	ev, ok := x.(company.UpdatedEvent)
	if !ok {
		return false
	}

	return m.expectedEv.CompanyID == ev.CompanyID && m.expectedEv.CompanyName == ev.CompanyName
}

func (m eventMatcher) String() string {
	return "match company updated event"
}

func TestUsecase_Execute_CreatesEvent(t *testing.T) {
	deps := setup(t)
	compID := uuid.New()
	oldName := "old name"
	newName := "new name"

	req := &update.Request{ID: compID, Name: &newName}

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	expectedEv := company.UpdatedEvent{CompanyID: compID, CompanyName: newName}

	deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, compID).
		Return(&member.CompanyMember{Role: member.CompanyRoleAdmin}, nil)

	deps.compRepo.EXPECT().UpdateAndGetOldName(gomock.Any(), req).
		Return(oldName, nil)

	deps.outbox.EXPECT().WriteCompanyUpdated(gomock.Any(), eventMatcher{expectedEv})

	deps.compCache.EXPECT().Del(gomock.Any(), compID)

	uc := NewUC(deps)

	err := uc.Execute(context.Background(), req, ident)

	require.NoError(t, err)
}

func TestUsecase_Execute_NoEvent(t *testing.T) {
	compID := uuid.New()
	oldName := "old name"

	req := &update.Request{ID: compID, Website: ptr("out new website")}

	tests := []struct {
		name string
		req  *update.Request
	}{
		{
			name: "no trigger data in request",
			req:  &update.Request{ID: compID, Website: ptr("out new website")},
		},
		{
			name: "same old data in request",
			req:  &update.Request{ID: compID, Name: &oldName},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

			deps := setup(t)

			deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, compID).
				Return(&member.CompanyMember{Role: member.CompanyRoleAdmin}, nil)

			deps.compRepo.EXPECT().UpdateAndGetOldName(gomock.Any(), req).
				Return(oldName, nil)

			deps.compCache.EXPECT().Del(gomock.Any(), compID)

			uc := NewUC(deps)

			err := uc.Execute(context.Background(), req, ident)

			require.NoError(t, err)
		})
	}
}

func TestUsecase_NotFound(t *testing.T) {
	deps := setup(t)
	compID := uuid.New()

	req := &update.Request{ID: compID, Website: ptr("out new website")}

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	deps.memRepo.EXPECT().Get(gomock.Any(), ident.UserID, compID).
		Return(&member.CompanyMember{Role: member.CompanyRoleAdmin}, nil)

	deps.compRepo.EXPECT().UpdateAndGetOldName(gomock.Any(), req).
		Return("", company.ErrCompanyNotFound)

	uc := NewUC(deps)

	err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, company.ErrCompanyNotFound)
}

func TestUsecase_ValidateErr(t *testing.T) {
	deps := setup(t)
	compID := uuid.New()

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	uc := NewUC(deps)

	tests := []struct {
		name string
		req  update.Request
		err  error
	}{
		{
			name: "empty name",
			req: update.Request{
				ID:   compID,
				Name: ptr(""),
			},
			err: company.ErrCompanyInvalidNameLen,
		},
		{
			name: "too long name",
			req: update.Request{
				ID:   compID,
				Name: ptr(string(make([]byte, company.MaxCompanyNameLen+1))),
			},
			err: company.ErrCompanyInvalidNameLen,
		},
		{
			name: "too long desc",
			req: update.Request{
				ID:          compID,
				Description: ptr(string(make([]byte, company.MaxCompanyDescriptionLen+1))),
			},
			err: company.ErrCompanyInvalidDescriptionLen,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.Execute(context.Background(), &tt.req, ident)

			require.ErrorIs(t, err, tt.err)
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
