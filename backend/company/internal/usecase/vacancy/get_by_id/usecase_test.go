package get_vacancy_test

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/get_by_id"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type repoMock struct {
	mock.Mock
}

func (m *repoMock) GetByID(ctx context.Context, id uuid.UUID, companyID uuid.UUID) (*domain.Vacancy, error) {
	args := m.Called(ctx, id, companyID)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Vacancy), args.Error(1)
	}
	return nil, args.Error(1)
}

type cacheMock struct {
	mock.Mock
}

func (m *cacheMock) Put(ctx context.Context, key uuid.UUID, val *domain.Vacancy, exp time.Duration) {
	m.Called(ctx, key, val, exp)
}

func (m *cacheMock) Get(ctx context.Context, id uuid.UUID) *domain.Vacancy {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Vacancy)
	}
	return nil
}

func TestUsecase_Execute_CacheHit(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)

	id := uuid.New()
	compID := uuid.New()

	cache.On("Get", mock.Anything, mock.Anything).
		Return(&domain.Vacancy{ID: id, CompanyID: compID, Title: "Title"}).Once()

	uc := get_vacancy.NewUsecase(repo, cache)

	resp, err := uc.Execute(context.Background(), id, compID)

	require.NoError(t, err)
	assert.Equal(t, resp.ID, id)
	assert.Equal(t, resp.Title, "Title")
	cache.AssertExpectations(t)
	repo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
}

func TestUsecase_Execute_CacheMiss(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)

	id := uuid.New()
	compID := uuid.New()

	cache.On("Get", mock.Anything, mock.Anything).
		Return(nil).Once()

	repo.On("GetByID", mock.Anything, mock.Anything, mock.Anything).
		Return(&domain.Vacancy{ID: id, CompanyID: compID, Title: "Title"}, nil).Once()

	cache.On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once()

	uc := get_vacancy.NewUsecase(repo, cache)

	resp, err := uc.Execute(context.Background(), id, compID)

	require.NoError(t, err)
	assert.Equal(t, resp.ID, id)
	assert.Equal(t, resp.Title, "Title")
	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestUsecase_Execute_RepoErr(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)

	id := uuid.New()
	compID := uuid.New()

	cache.On("Get", mock.Anything, mock.Anything).
		Return(nil).Once()

	repo.On("GetByID", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("err: i. e. not found ")).Once()

	uc := get_vacancy.NewUsecase(repo, cache)

	_, err := uc.Execute(context.Background(), id, compID)

	require.Error(t, err)
	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
	cache.AssertNotCalled(t, "Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}
