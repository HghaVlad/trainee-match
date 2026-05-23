package listsearch_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	list "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/listsearch"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/listsearch/mocks"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/views"
)

type testDeps struct {
	repo  *mocks.MockVacancyRepo
	cache *mocks.MockResponseCacheRepo
}

func setup(t *testing.T) *testDeps {
	ctrl := gomock.NewController(t)

	return &testDeps{
		repo:  mocks.NewMockVacancyRepo(ctrl),
		cache: mocks.NewMockResponseCacheRepo(ctrl),
	}
}

func newUC(deps *testDeps) *list.Usecase {
	return list.NewUsecase(
		deps.repo,
		deps.cache,
	)
}

func TestUsecase_Execute_Success_PublishedAt(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUC(deps)

	req := &list.Request{
		Order:        list.OrderPublishedAtDesc,
		Limit:        20,
		Requirements: &list.Requirements{},
	}

	vacancies := []views.PublishedVacSummary{
		{
			ID:             uuid.New(),
			CompanyID:      uuid.New(),
			CompanyName:    "Acme",
			Title:          "Go Developer",
			WorkFormat:     vacancy.WorkFormatRemote,
			EmploymentType: vacancy.EmploymentTypeInternship,
			IsPaid:         true,
			PublishedAt:    time.Now(),
		},
	}

	nextCursor := &list.PublishedAtCursor{
		PublishedAt: time.Now(),
		ID:          uuid.New(),
	}

	deps.cache.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(nil)

	deps.repo.EXPECT().
		ListPublishedSummaries(
			gomock.Any(),
			req.Requirements,
			list.OrderPublishedAtDesc,
			(*list.PublishedAtCursor)(nil),
			req.Limit,
		).
		Return(&list.SearchResult{
			Vacancies:  vacancies,
			NextCursor: nextCursor,
			HasNext:    true,
		}, nil)

	deps.cache.EXPECT().
		Put(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			time.Second*20,
		)

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Len(t, resp.Vacancies, 1)
	require.NotNil(t, resp.NextCursor)
	require.NotEmpty(t, *resp.NextCursor)
}

func TestUsecase_Execute_Success_FromCache(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUC(deps)

	req := &list.Request{
		Requirements: &list.Requirements{},
		Order:        list.OrderPublishedAtDesc,
		Limit:        20,
	}

	expectedResp := &list.Response{
		Vacancies: []views.PublishedVacSummary{
			{
				ID:          uuid.New(),
				Title:       "Cached vacancy",
				PublishedAt: time.Now(),
			},
		},
	}

	deps.cache.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(expectedResp)

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	require.Equal(t, expectedResp, resp)
}

func TestUsecase_Execute_RelevanceWithoutQuery_ChangesOrder(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUC(deps)

	req := &list.Request{
		Order:         list.OrderRelevance,
		EncodedCursor: "cursor",
		Limit:         20,
		Requirements:  &list.Requirements{},
	}

	deps.cache.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(nil)

	deps.repo.EXPECT().
		ListPublishedSummaries(
			gomock.Any(),
			req.Requirements,
			list.OrderPublishedAtDesc,
			(*list.PublishedAtCursor)(nil),
			req.Limit,
		).
		Return(&list.SearchResult{}, nil)

	deps.cache.EXPECT().
		Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Equal(t, list.OrderPublishedAtDesc, req.Order)
	require.Empty(t, req.EncodedCursor)
}

func TestUsecase_Execute_Success_Relevance(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUC(deps)

	query := "golang"

	req := &list.Request{
		Order: list.OrderRelevance,
		Limit: 20,
		Requirements: &list.Requirements{
			Query: &query,
		},
	}

	deps.cache.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(nil)

	deps.repo.EXPECT().
		ListPublishedSummaries(
			gomock.Any(),
			req.Requirements,
			list.OrderRelevance,
			(*list.RelevanceCursor)(nil),
			req.Limit,
		).
		Return(&list.SearchResult{}, nil)

	deps.cache.EXPECT().
		Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestUsecase_Execute_Success_SalaryOrder(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUC(deps)

	isPaid := true

	req := &list.Request{
		Order: list.OrderSalaryDesc,
		Limit: 20,
		Requirements: &list.Requirements{
			IsPaid: &isPaid,
		},
	}

	deps.cache.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(nil)

	deps.repo.EXPECT().
		ListPublishedSummaries(
			gomock.Any(),
			req.Requirements,
			list.OrderSalaryDesc,
			(*list.SalaryCursor)(nil),
			req.Limit,
		).
		Return(&list.SearchResult{}, nil)

	deps.cache.EXPECT().
		Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestUsecase_Execute_InvalidRequest(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUC(deps)

	req := &list.Request{
		Order: "invalid",
	}

	resp, err := uc.Execute(context.Background(), req)

	require.ErrorIs(t, err, common.ErrUnsupportedListOrder)
	require.Nil(t, resp)
}

func TestUsecase_Execute_RepoError(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUC(deps)

	req := &list.Request{
		Requirements: &list.Requirements{},
		Order:        list.OrderPublishedAtDesc,
		Limit:        20,
	}

	expectedErr := errors.New("repo error")

	deps.cache.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(nil)

	deps.repo.EXPECT().
		ListPublishedSummaries(
			gomock.Any(),
			gomock.Any(),
			list.OrderPublishedAtDesc,
			(*list.PublishedAtCursor)(nil),
			req.Limit,
		).
		Return(nil, expectedErr)

	resp, err := uc.Execute(context.Background(), req)

	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, resp)
}

func TestUsecase_Execute_InvalidCursor(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUC(deps)

	req := &list.Request{
		Requirements:  new(list.Requirements),
		Order:         list.OrderPublishedAtDesc,
		Limit:         20,
		EncodedCursor: "invalid-cursor",
	}

	deps.cache.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(nil)

	resp, err := uc.Execute(context.Background(), req)

	require.Error(t, err)
	require.Nil(t, resp)
}

func TestUsecase_Execute_HasNoNext(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUC(deps)

	req := &list.Request{
		Requirements: new(list.Requirements),
		Order:        list.OrderPublishedAtDesc,
		Limit:        20,
	}

	deps.cache.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(nil)

	deps.repo.EXPECT().
		ListPublishedSummaries(
			gomock.Any(),
			gomock.Any(),
			list.OrderPublishedAtDesc,
			(*list.PublishedAtCursor)(nil),
			req.Limit,
		).
		Return(&list.SearchResult{
			Vacancies: []views.PublishedVacSummary{},
			HasNext:   false,
		}, nil)

	deps.cache.EXPECT().
		Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Nil(t, resp.NextCursor)
}
