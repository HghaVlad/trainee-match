package list_companies_test

import (
	"context"
	"testing"
	"time"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/list"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type companyRepoMock struct {
	mock.Mock
}

func (m *companyRepoMock) ListByCreatedAtDesc(
	ctx context.Context,
	cursor *list_companies.CreatedAtCursor,
	limit int,
) ([]list_companies.CompanySummary, *list_companies.CreatedAtCursor, error) {

	args := m.Called(ctx, cursor, limit)

	cs := args.Get(0)
	next := args.Get(1)

	if cs != nil && next != nil {
		return cs.([]list_companies.CompanySummary), next.(*list_companies.CreatedAtCursor), args.Error(2)
	}

	if cs != nil {
		return cs.([]list_companies.CompanySummary), nil, args.Error(2)
	}

	return nil, nil, args.Error(2)
}

func (m *companyRepoMock) ListByName(
	ctx context.Context,
	cursor *list_companies.NameCursor,
	limit int,
) ([]list_companies.CompanySummary, *list_companies.NameCursor, error) {

	args := m.Called(ctx, cursor, limit)

	cs := args.Get(0)
	next := args.Get(1)

	if cs != nil && next != nil {
		return cs.([]list_companies.CompanySummary), next.(*list_companies.NameCursor), args.Error(2)
	}

	if cs != nil {
		return cs.([]list_companies.CompanySummary), nil, args.Error(2)
	}

	return nil, nil, args.Error(2)
}

func (m *companyRepoMock) ListByVacanciesCnt(
	ctx context.Context,
	cursor *list_companies.VacanciesCntCursor,
	limit int,
) ([]list_companies.CompanySummary, *list_companies.VacanciesCntCursor, error) {

	args := m.Called(ctx, cursor, limit)

	cs := args.Get(0)
	next := args.Get(1)

	if cs != nil && next != nil {
		return cs.([]list_companies.CompanySummary), next.(*list_companies.VacanciesCntCursor), args.Error(2)
	}

	if cs != nil {
		return cs.([]list_companies.CompanySummary), nil, args.Error(2)
	}

	return nil, nil, args.Error(2)
}

type responseCacheRepoMock struct {
	mock.Mock
}

func (m *responseCacheRepoMock) Get(ctx context.Context, key string) *list_companies.Response {
	args := m.Called(ctx, key)

	if resp := args.Get(0); resp != nil {
		return resp.(*list_companies.Response)
	}

	return nil
}

func (m *responseCacheRepoMock) Put(ctx context.Context, key string, response *list_companies.Response, exp time.Duration) {
	m.Called(ctx, key, response, exp)
}

func TestExecute_CacheHit(t *testing.T) {
	repo := new(companyRepoMock)
	cache := new(responseCacheRepoMock)

	req := &list_companies.Request{
		Order:  list_companies.OrderVacanciesDesc,
		Cursor: "curs",
		Limit:  10,
	}

	expectedResp := &list_companies.Response{}

	cache.
		On("Get", mock.Anything, mock.Anything).
		Return(expectedResp).Once()

	uc := list_companies.NewUsecase(repo, cache)

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	require.Equal(t, expectedResp, resp)

	cache.AssertExpectations(t)
	repo.AssertNotCalled(t, "ListByVacanciesCnt", mock.Anything)
}

func TestExecute_CacheMiss(t *testing.T) {
	repo := new(companyRepoMock)
	cache := new(responseCacheRepoMock)

	req := &list_companies.Request{
		Order:  list_companies.OrderVacanciesDesc,
		Cursor: "",
		Limit:  10,
	}

	cache.
		On("Get", mock.Anything, mock.Anything).
		Return(nil).Once()

	repo.
		On("ListByVacanciesCnt", mock.Anything, mock.Anything, mock.Anything).
		Return([]list_companies.CompanySummary{}, nil, nil).Once()

	cache.
		On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once()

	uc := list_companies.NewUsecase(repo, cache)

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	require.NotEmpty(t, resp)

	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestListByCreatedAt_OK(t *testing.T) {
	req := &list_companies.Request{
		Order:  list_companies.OrderCreatedAtDesc,
		Cursor: "",
		Limit:  10,
	}

	comps := []list_companies.CompanySummary{{}}

	t.Run("ok: has next page", func(t *testing.T) {
		repo := new(companyRepoMock)
		uc := list_companies.NewUsecase(repo, nil)

		nextCursor := list_companies.CreatedAtCursor{
			CreatedAt: time.Now(),
			Name:      "Acme",
		}

		repo.
			On("ListByCreatedAtDesc", mock.Anything, mock.Anything, 10).
			Return(comps, &nextCursor, nil)

		resp, err := uc.ListByCreatedAt(context.Background(), req)

		require.NoError(t, err)
		assert.Len(t, resp.Companies, 1)
		assert.NotEmpty(t, resp.NextCursor)
	})

	t.Run("ok: no next page", func(t *testing.T) {
		repo := new(companyRepoMock)
		uc := list_companies.NewUsecase(repo, nil)

		repo.
			On("ListByCreatedAtDesc", mock.Anything, mock.Anything, 10).
			Return(comps, nil, nil)

		resp, err := uc.ListByCreatedAt(context.Background(), req)

		require.NoError(t, err)
		assert.Len(t, resp.Companies, 1)
		assert.Empty(t, resp.NextCursor)
	})
}
