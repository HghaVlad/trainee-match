package delete_company_test

import (
	"context"
	"errors"
	"testing"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/delete"
)

type repoMock struct {
	mock.Mock
}

func (r *repoMock) Delete(ctx context.Context, id uuid.UUID) error {
	return r.Called(ctx, id).Error(0)
}

type cacheMock struct {
	mock.Mock
}

func (r *cacheMock) Del(ctx context.Context, id uuid.UUID) {
	r.Called(ctx, id)
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

func TestUsecase_Execute_OK(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)
	memRepo := new(memRepoMock)

	repo.On("Delete", mock.Anything, mock.Anything).
		Return(nil).Once()

	memRepo.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(&domain.CompanyMember{Role: value_types.CompanyRoleAdmin}, nil).Once()

	cache.On("Del", mock.Anything, mock.Anything).Once()

	uc := delete_company.NewUsecase(repo, memRepo, cache)

	idenity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

	err := uc.Execute(context.Background(), uuid.New(), idenity)

	require.NoError(t, err)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestUsecase_Execute_CompanyRepoErr(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)
	memRepo := new(memRepoMock)

	memRepo.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(&domain.CompanyMember{Role: value_types.CompanyRoleAdmin}, nil).Once()

	repo.On("Delete", mock.Anything, mock.Anything).
		Return(errors.New("some domain err")).Once()

	uc := delete_company.NewUsecase(repo, memRepo, cache)

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

	err := uc.Execute(context.Background(), uuid.New(), identity)

	require.Error(t, err)
	repo.AssertExpectations(t)
	cache.AssertNotCalled(t, "Del", mock.Anything, mock.Anything)
}
