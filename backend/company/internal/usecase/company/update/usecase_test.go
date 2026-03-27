package update_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/update"
)

type companyRepoMock struct {
	mock.Mock
}

func (m *companyRepoMock) Update(ctx context.Context, req *update.Request) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

type cacheMock struct {
	mock.Mock
}

func (m *cacheMock) Del(ctx context.Context, id uuid.UUID) {
	m.Called(ctx, id)
}

type memRepoMock struct {
	mock.Mock
}

func (m *memRepoMock) Get(ctx context.Context, userID, companyID uuid.UUID) (*member.CompanyMember, error) {
	res := m.Called(ctx, userID, companyID)

	if c := res.Get(0); c != nil {
		return c.(*member.CompanyMember), res.Error(1)
	}

	return nil, res.Error(1)
}

func TestUsecase_ExecuteOK(t *testing.T) {
	cache := new(cacheMock)
	memRepo := new(memRepoMock)
	repo := new(companyRepoMock)

	memRepo.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(&member.CompanyMember{Role: member.CompanyRoleAdmin}, nil).Once()

	repo.On("Update", mock.Anything, mock.Anything).
		Return(nil).Once()

	cache.On("Del", mock.Anything, mock.Anything).Once()

	uc := update.NewUsecase(repo, memRepo, cache)

	req := &update.Request{
		ID:   uuid.New(),
		Name: ptr("Acme"),
	}

	ident := identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	err := uc.Execute(context.Background(), req, ident)

	require.NoError(t, err)
	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestUsecase_DbErr(t *testing.T) {
	cache := new(cacheMock)
	memRepo := new(memRepoMock)
	repo := new(companyRepoMock)

	memRepo.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(&member.CompanyMember{Role: member.CompanyRoleAdmin}, nil).Once()

	repo.On("Update", mock.Anything, mock.Anything).
		Return(errors.New("db err")).Once()

	uc := update.NewUsecase(repo, memRepo, cache)

	req := &update.Request{
		ID:   uuid.New(),
		Name: ptr("Acme"),
	}

	ident := identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	err := uc.Execute(context.Background(), req, ident)

	assert.Error(t, err)
	repo.AssertExpectations(t)
	cache.AssertNotCalled(t, "Del", mock.Anything, mock.Anything)
}

func TestUsecase_ValidateErr(t *testing.T) {
	cache := new(cacheMock)
	repo := new(companyRepoMock)
	memRepo := new(memRepoMock)

	ident := identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	uc := update.NewUsecase(repo, memRepo, cache)

	tests := []struct {
		name string
		req  update.Request
		err  error
	}{
		{
			name: "empty name",
			req: update.Request{
				ID:   uuid.New(),
				Name: ptr(""),
			},
			err: company.ErrCompanyInvalidNameLen,
		},
		{
			name: "too long name",
			req: update.Request{
				ID:   uuid.New(),
				Name: ptr(string(make([]byte, company.MaxCompanyNameLen+1))),
			},
			err: company.ErrCompanyInvalidNameLen,
		},
		{
			name: "too long desc",
			req: update.Request{
				ID:          uuid.New(),
				Description: ptr(string(make([]byte, company.MaxCompanyDescriptionLen+1))),
			},
			err: company.ErrCompanyInvalidDescriptionLen,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.Execute(context.Background(), &tt.req, ident)

			assert.Equal(t, tt.err, err)
			repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			cache.AssertNotCalled(t, "Del", mock.Anything, mock.Anything)
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
