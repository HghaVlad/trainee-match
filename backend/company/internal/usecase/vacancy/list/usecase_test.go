package list_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list"
)

type cacheMock struct {
	mock.Mock
}

func (m *cacheMock) Get(ctx context.Context, key string) *list.Response {
	args := m.Called(ctx, key)
	if args.Get(0) != nil {
		return args.Get(0).(*list.Response)
	}
	return nil
}

func (m *cacheMock) Put(ctx context.Context, key string, response *list.Response, exp time.Duration) {
	m.Called(ctx, key, response, exp)
}

type repoMock struct {
	mock.Mock
}

func (m *repoMock) ListPublishedSummaries(
	ctx context.Context,
	requirements *list.Requirements,
	order list.Order,
	cursor any,
	limit int,
) ([]list.VacancySummary, error) {
	args := m.Called(ctx, requirements, order, cursor, limit)

	vcs := args.Get(0)

	if vcs != nil {
		return vcs.([]list.VacancySummary), args.Error(1)
	}

	return nil, args.Error(1)
}

func TestUsecase_Execute_CacheHit(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)

	req := &list.Request{
		Order:         list.OrderPublishedAtDesc,
		EncodedCursor: "",
		Limit:         10,
	}

	cache.On("Get", mock.Anything, mock.Anything).
		Return(&list.Response{Vacancies: []list.VacancySummary{{}}}).Once()

	uc := list.NewUsecase(repo, cache)

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	assert.Len(t, resp.Vacancies, 1)
	cache.AssertExpectations(t)
	repo.AssertNotCalled(
		t,
		"ListPublishedSummaries",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	)
}

func TestUsecase_Execute_NextCursor(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)

	req := &list.Request{
		Order:         list.OrderPublishedAtDesc,
		EncodedCursor: "",
		Limit:         10,
	}

	vcs := make([]list.VacancySummary, req.Limit+1)
	for i := range vcs {
		vcs[i] = list.VacancySummary{
			ID:          uuid.New(),
			PublishedAt: time.Now().Add(-time.Duration(i) * time.Minute),
		}
	}

	cache.On("Get", mock.Anything, mock.Anything).
		Return(nil).Once()

	repo.On("ListPublishedSummaries", mock.Anything, mock.Anything, mock.Anything, mock.Anything, req.Limit+1).
		Return(vcs, nil).Once()

	cache.On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once()

	uc := list.NewUsecase(repo, cache)

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

	req := &list.Request{
		Order:         list.OrderPublishedAtDesc,
		EncodedCursor: "",
		Limit:         10,
	}

	vcs := make([]list.VacancySummary, req.Limit)
	for i := range vcs {
		vcs[i] = list.VacancySummary{
			ID:          uuid.New(),
			PublishedAt: time.Now().Add(-time.Duration(i) * time.Minute),
		}
	}

	cache.On("Get", mock.Anything, mock.Anything).
		Return(nil).Once()

	repo.On("ListPublishedSummaries", mock.Anything, mock.Anything, mock.Anything, mock.Anything, req.Limit+1).
		Return(vcs, nil).Once()

	cache.On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once()

	uc := list.NewUsecase(repo, cache)

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, len(resp.Vacancies), req.Limit)
	assert.Empty(t, resp.NextCursor)
	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}
