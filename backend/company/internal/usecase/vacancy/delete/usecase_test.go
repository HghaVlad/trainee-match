package delete_vacancy_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
	uc_common "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/delete"
)

type repoMock struct {
	mock.Mock
}

func (m *repoMock) Delete(ctx context.Context, vacancyID uuid.UUID, companyID uuid.UUID) error {
	return m.Called(ctx, vacancyID, companyID).Error(0)
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

func TestUsecase_Execute_HappyPath(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)
	memRepo := new(memRepoMock)

	memRepo.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(&domain.CompanyMember{Role: value_types.CompanyRoleAdmin}, nil).Once()

	repo.On("Delete", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()

	cache.On("Del", mock.Anything, mock.Anything).
		Return(nil).Once()

	uc := delete_vacancy.NewUsecase(repo, memRepo, cache)

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

	err := uc.Execute(context.Background(), uuid.New(), uuid.New(), identity)

	require.NoError(t, err)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}
