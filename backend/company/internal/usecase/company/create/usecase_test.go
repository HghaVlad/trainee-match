package create_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/create"
)

type repoMock struct {
	mock.Mock
}

func (m *repoMock) Create(ctx context.Context, company *company.Company) error {
	return m.Called(ctx, company).Error(0)
}

type companyMemberRepoMock struct {
	mock.Mock
}

func (m *companyMemberRepoMock) Create(ctx context.Context, member *member.CompanyMember) error {
	return m.Called(ctx, member).Error(0)
}

type FakeTxManager struct {
	Called bool
}

func (f *FakeTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	f.Called = true
	return fn(ctx)
}

func TestUsecase_ExecuteOK(t *testing.T) {
	repo := new(repoMock)
	memRepo := new(companyMemberRepoMock)
	txManager := new(FakeTxManager)

	req := &create.Request{
		Name:        "Acme",
		Description: ptr("Hello!"),
	}

	repo.On("Create", mock.Anything, mock.Anything).
		Return(nil).Once()

	memRepo.On("Create", mock.Anything, mock.Anything).
		Return(nil).Once()

	uc := create.NewUsecase(repo, memRepo, txManager)

	iden := identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleHR,
	}

	resp, err := uc.Execute(context.Background(), req, iden)

	require.NoError(t, err)
	require.NotNil(t, resp)
	repo.AssertExpectations(t)
}

func TestUsecase_ExecuteValidateErr(t *testing.T) {
	repo := new(repoMock)
	memRepo := new(companyMemberRepoMock)
	txManager := new(FakeTxManager)

	iden := identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleHR,
	}

	uc := create.NewUsecase(repo, memRepo, txManager)

	tests := []struct {
		name string
		req  create.Request
		err  error
	}{
		{
			name: "empty name",
			req: create.Request{
				Name: "",
			},
			err: company.ErrCompanyInvalidNameLen,
		},
		{
			name: "too long name",
			req: create.Request{
				Name: string(make([]byte, company.MaxCompanyNameLen+1)),
			},
			err: company.ErrCompanyInvalidNameLen,
		},
		{
			name: "too long desc",
			req: create.Request{
				Name:        "Acme",
				Description: ptr(string(make([]byte, company.MaxCompanyDescriptionLen+1))),
			},
			err: company.ErrCompanyInvalidDescriptionLen,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), &tt.req, iden)

			assert.Equal(t, tt.err, err)
			repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
