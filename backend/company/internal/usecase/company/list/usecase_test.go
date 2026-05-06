package list_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/list"
)

type companyRepoMock struct {
	mock.Mock
}

func (c *companyRepoMock) ListSummaries(
	ctx context.Context,
	order list.Order,
	filter list.Filter,
	cursor any,
	limit int,
) ([]list.CompanySummary, error) {
	args := c.Called(ctx, order, filter, cursor, limit)

	if comps := args.Get(0); comps != nil {
		return comps.([]list.CompanySummary), nil
	}

	return nil, args.Error(1)
}

type responseCacheRepoMock struct {
	mock.Mock
}

func (m *responseCacheRepoMock) Get(ctx context.Context, key string) *list.Response {
	args := m.Called(ctx, key)

	if resp := args.Get(0); resp != nil {
		return resp.(*list.Response)
	}

	return nil
}

func (m *responseCacheRepoMock) Put(ctx context.Context, key string, response *list.Response, exp time.Duration) {
	m.Called(ctx, key, response, exp)
}

func TestExecute_CacheHit(t *testing.T) {
	repo := new(companyRepoMock)
	cache := new(responseCacheRepoMock)

	req := &list.Request{
		Order:         list.OrderVacanciesDesc,
		EncodedCursor: "curs",
		Limit:         10,
	}

	expectedResp := &list.Response{}

	cache.
		On("Get", mock.Anything, mock.Anything).
		Return(expectedResp).Once()

	uc := list.NewUsecase(repo, cache)

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	require.Equal(t, expectedResp, resp)

	cache.AssertExpectations(t)
	repo.AssertNotCalled(t, "ListSummaries", mock.Anything)
}

func TestExecute_CacheMiss(t *testing.T) {
	repo := new(companyRepoMock)
	cache := new(responseCacheRepoMock)

	req := &list.Request{
		Order:         list.OrderVacanciesDesc,
		EncodedCursor: "",
		Limit:         10,
	}

	cache.
		On("Get", mock.Anything, mock.Anything).
		Return(nil).Once()

	repo.
		On("ListSummaries", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]list.CompanySummary{}, nil).Once()

	cache.
		On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once()

	uc := list.NewUsecase(repo, cache)

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	require.NotEmpty(t, resp)

	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}
