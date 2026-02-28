package list_vacancy_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list"
)

type cacheMock struct {
	mock.Mock
}

func (m *cacheMock) Get(ctx context.Context, key string) *list_vacancy.Response {
	args := m.Called(ctx, key)
	if args.Get(0) != nil {
		return args.Get(0).(*list_vacancy.Response)
	}
	return nil
}

func (m *cacheMock) Put(ctx context.Context, key string, response *list_vacancy.Response, exp time.Duration) {
	m.Called(ctx, key, response, exp)
}

type repoMock struct {
	mock.Mock
}

func (m *repoMock) List(ctx context.Context, requirements *list_vacancy.Requirements, order list_vacancy.Order, cursor any, limit int) ([]list_vacancy.VacancySummary, error) {
	args := m.Called(ctx, requirements, order, cursor, limit)

	vcs := args.Get(0)

	if vcs != nil {
		return vcs.([]list_vacancy.VacancySummary), args.Error(1)
	}

	return nil, args.Error(1)
}

func TestUsecase_Execute_CacheHit(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)

	req := &list_vacancy.Request{
		Order:         list_vacancy.OrderPublishedAtDesc,
		EncodedCursor: "",
		Limit:         10,
	}

	cache.On("Get", mock.Anything, mock.Anything).
		Return(&list_vacancy.Response{Vacancies: []list_vacancy.VacancySummary{{}}}).Once()

	uc := list_vacancy.NewUsecase(repo, cache)

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, len(resp.Vacancies), 1)
	cache.AssertExpectations(t)
	repo.AssertNotCalled(t, "ListByPublishedAt", mock.Anything, mock.Anything, mock.Anything)
}

func TestUsecase_Execute_NextCursor(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)

	req := &list_vacancy.Request{
		Order:         list_vacancy.OrderPublishedAtDesc,
		EncodedCursor: "",
		Limit:         10,
	}

	vcs := make([]list_vacancy.VacancySummary, req.Limit+1)

	cache.On("Get", mock.Anything, mock.Anything).
		Return(nil).Once()

	repo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, req.Limit+1).
		Return(vcs, nil).Once()

	cache.On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once()

	uc := list_vacancy.NewUsecase(repo, cache)

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, len(resp.Vacancies), req.Limit)
	assert.NotEmpty(t, resp.NextCursor)
	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestUsecase_Execute_NoNextCursor(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)

	req := &list_vacancy.Request{
		Order:         list_vacancy.OrderPublishedAtDesc,
		EncodedCursor: "",
		Limit:         10,
	}

	vcs := make([]list_vacancy.VacancySummary, req.Limit)

	cache.On("Get", mock.Anything, mock.Anything).
		Return(nil).Once()

	repo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, req.Limit+1).
		Return(vcs, nil).Once()

	cache.On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once()

	uc := list_vacancy.NewUsecase(repo, cache)

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, len(resp.Vacancies), req.Limit)
	assert.Empty(t, resp.NextCursor)
	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}
