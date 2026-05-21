package listsearch_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/listsearch"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/views"
)

type cacheMock struct {
	mock.Mock
}

func (m *cacheMock) Get(ctx context.Context, key string) *listsearch.Response {
	args := m.Called(ctx, key)
	if args.Get(0) != nil {
		return args.Get(0).(*listsearch.Response)
	}
	return nil
}

func (m *cacheMock) Put(ctx context.Context, key string, response *listsearch.Response, exp time.Duration) {
	m.Called(ctx, key, response, exp)
}

type repoMock struct {
	mock.Mock
}

func (m *repoMock) ListPublishedSummaries(
	ctx context.Context,
	requirements *listsearch.Requirements,
	order listsearch.Order,
	cursor any,
	limit int,
) ([]views.PublishedVacSummary, error) {
	args := m.Called(ctx, requirements, order, cursor, limit)

	vcs := args.Get(0)

	if vcs != nil {
		return vcs.([]views.PublishedVacSummary), args.Error(1)
	}

	return nil, args.Error(1)
}

func TestUsecase_Execute_CacheHit(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)

	req := &listsearch.Request{
		Order:         listsearch.OrderPublishedAtDesc,
		EncodedCursor: "",
		Limit:         10,
	}

	cache.On("Get", mock.Anything, mock.Anything).
		Return(&listsearch.Response{Vacancies: []views.PublishedVacSummary{{}}}).Once()

	uc := listsearch.NewUsecase(repo, cache)

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

	req := &listsearch.Request{
		Order:         listsearch.OrderPublishedAtDesc,
		EncodedCursor: "",
		Limit:         10,
	}

	vcs := make([]views.PublishedVacSummary, req.Limit+1)
	for i := range vcs {
		vcs[i] = views.PublishedVacSummary{
			ID:          uuid.New(),
			PublishedAt: time.Now().Add(-time.Duration(i) * time.Minute),
		}
	}

	cache.On("Get", mock.Anything, mock.Anything).
		Return(nil).Once()

	repo.On("ListPublishedSummaries", mock.Anything, mock.Anything, mock.Anything, mock.Anything, req.Limit+1).
		Return(vcs, nil).Once()

	cache.On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once()

	uc := listsearch.NewUsecase(repo, cache)

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

	req := &listsearch.Request{
		Order:         listsearch.OrderPublishedAtDesc,
		EncodedCursor: "",
		Limit:         10,
	}

	vcs := make([]views.PublishedVacSummary, req.Limit)
	for i := range vcs {
		vcs[i] = views.PublishedVacSummary{
			ID:          uuid.New(),
			PublishedAt: time.Now().Add(-time.Duration(i) * time.Minute),
		}
	}

	cache.On("Get", mock.Anything, mock.Anything).
		Return(nil).Once()

	repo.On("ListPublishedSummaries", mock.Anything, mock.Anything, mock.Anything, mock.Anything, req.Limit+1).
		Return(vcs, nil).Once()

	cache.On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once()

	uc := listsearch.NewUsecase(repo, cache)

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, len(resp.Vacancies), req.Limit)
	assert.Empty(t, resp.NextCursor)
	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}
