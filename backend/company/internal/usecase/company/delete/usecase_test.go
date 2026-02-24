package delete_company_test

import (
	"context"
	"errors"
	"testing"

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

func TestUsecase_Execute_OK(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)

	repo.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()

	cache.On("Del", mock.Anything, mock.Anything).Once()

	uc := delete_company.NewUsecase(repo, cache)

	err := uc.Execute(context.Background(), uuid.New())

	require.NoError(t, err)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestUsecase_Execute_RepoErr(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)

	repo.On("Delete", mock.Anything, mock.Anything).
		Return(errors.New("some domain err")).Once()

	uc := delete_company.NewUsecase(repo, cache)

	err := uc.Execute(context.Background(), uuid.New())

	require.Error(t, err)
	repo.AssertExpectations(t)
	cache.AssertNotCalled(t, "Del", mock.Anything, mock.Anything)
}
