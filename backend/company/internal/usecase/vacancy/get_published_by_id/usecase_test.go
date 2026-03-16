package get_published_vacancy_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/get_published_by_id"
)

type repoMock struct {
	mock.Mock
}

func (m *repoMock) GetPublishedByID(ctx context.Context, id uuid.UUID) (*get_published_vacancy.Response, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*get_published_vacancy.Response), args.Error(1)
	}
	return nil, args.Error(1)
}

type cacheMock struct {
	mock.Mock
}

func (m *cacheMock) Put(ctx context.Context, key uuid.UUID, val *get_published_vacancy.Response, exp time.Duration) {
	m.Called(ctx, key, val, exp)
}

func (m *cacheMock) Get(ctx context.Context, id uuid.UUID) *get_published_vacancy.Response {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*get_published_vacancy.Response)
	}
	return nil
}

func TestUsecase_Execute_CacheHit(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)

	id := uuid.New()

	cache.On("Get", mock.Anything, mock.Anything).
		Return(&get_published_vacancy.Response{ID: id, Title: "Title"}).Once()

	uc := get_published_vacancy.NewUsecase(repo, cache)

	resp, err := uc.Execute(context.Background(), id)

	require.NoError(t, err)
	assert.Equal(t, resp.ID, id)
	assert.Equal(t, resp.Title, "Title")
	cache.AssertExpectations(t)
	repo.AssertNotCalled(t, "GetPublishedByID", mock.Anything, mock.Anything)
}

func TestUsecase_Execute_CacheMiss(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)

	id := uuid.New()

	cache.On("Get", mock.Anything, mock.Anything).
		Return(nil).Once()

	repo.On("GetPublishedByID", mock.Anything, mock.Anything).
		Return(&get_published_vacancy.Response{ID: id, Title: "Title"}, nil).Once()

	cache.On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once()

	uc := get_published_vacancy.NewUsecase(repo, cache)

	resp, err := uc.Execute(context.Background(), id)

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

	cache.On("Get", mock.Anything, mock.Anything).
		Return(nil).Once()

	repo.On("GetPublishedByID", mock.Anything, mock.Anything).
		Return(nil, errors.New("err: i. e. not found ")).Once()

	uc := get_published_vacancy.NewUsecase(repo, cache)

	_, err := uc.Execute(context.Background(), id)

	require.Error(t, err)
	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
	cache.AssertNotCalled(t, "Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}
