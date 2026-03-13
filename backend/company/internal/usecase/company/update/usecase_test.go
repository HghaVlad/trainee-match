package update_company_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/update"
)

type companyRepoMock struct {
	mock.Mock
}

func (m *companyRepoMock) Update(ctx context.Context, req *update_company.Request) error {
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

func (m *memRepoMock) Get(ctx context.Context, userID, companyID uuid.UUID) (*domain.CompanyMember, error) {
	res := m.Called(ctx, userID, companyID)

	if c := res.Get(0); c != nil {
		return c.(*domain.CompanyMember), res.Error(1)
	}

	return nil, res.Error(1)
}

func TestUsecase_ExecuteOK(t *testing.T) {
	cache := new(cacheMock)
	memRepo := new(memRepoMock)
	repo := new(companyRepoMock)

	memRepo.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(&domain.CompanyMember{Role: value_types.CompanyRoleAdmin}, nil).Once()

	repo.On("Update", mock.Anything, mock.Anything).
		Return(nil).Once()

	cache.On("Del", mock.Anything, mock.Anything).Once()

	uc := update_company.NewUsecase(repo, memRepo, cache)

	req := &update_company.Request{
		ID:   uuid.New(),
		Name: ptr("Acme"),
	}

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

	err := uc.Execute(context.Background(), req, identity)

	require.NoError(t, err)
	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestUsecase_DbErr(t *testing.T) {
	cache := new(cacheMock)
	memRepo := new(memRepoMock)
	repo := new(companyRepoMock)

	memRepo.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(&domain.CompanyMember{Role: value_types.CompanyRoleAdmin}, nil).Once()

	repo.On("Update", mock.Anything, mock.Anything).
		Return(errors.New("db err")).Once()

	uc := update_company.NewUsecase(repo, memRepo, cache)

	req := &update_company.Request{
		ID:   uuid.New(),
		Name: ptr("Acme"),
	}

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

	err := uc.Execute(context.Background(), req, identity)

	assert.Error(t, err)
	repo.AssertExpectations(t)
	cache.AssertNotCalled(t, "Del", mock.Anything, mock.Anything)
}

func TestUsecase_ValidateErr(t *testing.T) {
	cache := new(cacheMock)
	repo := new(companyRepoMock)
	memRepo := new(memRepoMock)

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

	uc := update_company.NewUsecase(repo, memRepo, cache)

	tests := []struct {
		name string
		req  update_company.Request
		err  error
	}{
		{
			name: "empty name",
			req: update_company.Request{
				ID:   uuid.New(),
				Name: ptr(""),
			},
			err: domain_errors.ErrCompanyInvalidNameLen,
		},
		{
			name: "too long name",
			req: update_company.Request{
				ID:   uuid.New(),
				Name: ptr(string(make([]byte, domain.MaxCompanyNameLen+1))),
			},
			err: domain_errors.ErrCompanyInvalidNameLen,
		},
		{
			name: "too long desc",
			req: update_company.Request{
				ID:          uuid.New(),
				Description: ptr(string(make([]byte, domain.MaxCompanyDescriptionLen+1))),
			},
			err: domain_errors.ErrCompanyInvalidDescriptionLen,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.Execute(context.Background(), &tt.req, identity)

			assert.Equal(t, tt.err, err)
			repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			cache.AssertNotCalled(t, "Del", mock.Anything, mock.Anything)
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
