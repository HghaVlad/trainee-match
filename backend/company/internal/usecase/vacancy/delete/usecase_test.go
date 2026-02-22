package delete_vacancy_test

import (
	"context"
	"testing"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/delete"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

func TestUsecase_Execute_HappyPath(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)

	repo.On("Delete", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()

	cache.On("Del", mock.Anything, mock.Anything).
		Return(nil).Once()

	uc := delete_vacancy.NewUsecase(repo, cache)

	err := uc.Execute(context.Background(), uuid.New(), uuid.New())

	require.NoError(t, err)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}
