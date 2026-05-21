package listsearch

import (
	"context"
	"time"

	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/utils/encoding"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/views"
)

// Usecase of vacancy listing, uses cursor pagination.
// Supports order by published_at, salary.
// Supports filters in Requirements.
type Usecase struct {
	repo      VacancyRepo
	respCache ResponseCacheRepo
}

func NewUsecase(repo VacancyRepo, cache ResponseCacheRepo) *Usecase {
	return &Usecase{repo: repo, respCache: cache}
}

// Execute cursor pagination list vacancy.
// Supports order by published_at, salary.
// Supports filters in Requirements.
func (uc *Usecase) Execute(ctx context.Context, req *Request) (*Response, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	respCacheKey := requestToCacheKey(req)
	resp := uc.respCache.Get(ctx, respCacheKey)
	if resp != nil {
		return resp, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	var err error

	switch req.Order {
	case OrderRelevance:
		resp, err = list[RelevanceCursor](ctx, uc, req)
	case OrderPublishedAtDesc:
		resp, err = list[PublishedAtCursor](ctx, uc, req)
	case OrderSalaryDesc, OrderSalaryAsc:
		resp, err = list[SalaryCursor](ctx, uc, req)

	default:
		return nil, common.ErrUnsupportedListOrder
	}

	if err != nil {
		return nil, err
	}

	// Adding to cache with short ttl because it won't be updated/deleted by service
	uc.respCache.Put(ctx, respCacheKey, resp, time.Second*20)

	return resp, nil
}

func list[CursorT any](ctx context.Context, uc *Usecase, req *Request) (*Response, error) {
	cursor, curErr := encoding.DecodeCursor[CursorT, Order](req.EncodedCursor, req.Order)
	if curErr != nil {
		return nil, curErr
	}

	res, err := uc.repo.ListPublishedSummaries(ctx, req.Requirements, req.Order, cursor, req.Limit)
	if err != nil {
		return nil, err
	}

	if res.HasNext {
		nextCursor, _ := res.NextCursor.(*CursorT)
		return buildResponse[CursorT](res.Vacancies, nextCursor, req.Order)
	}

	return buildResponse[CursorT](res.Vacancies, nil, req.Order)
}

func getNextCursor[CursorT any](
	vacancies []views.PublishedVacSummary,
	limit int,
) (*CursorT, []views.PublishedVacSummary) {
	if len(vacancies) <= limit {
		return nil, vacancies
	}

	vacancies = vacancies[:len(vacancies)-1]
	last := vacancies[len(vacancies)-1]

	var zero CursorT
	var cursor any

	switch any(zero).(type) {
	case PublishedAtCursor:
		cursor = &PublishedAtCursor{
			PublishedAt: last.PublishedAt,
			ID:          last.ID,
		}

	case SalaryCursor:
		if last.SalaryFrom == nil || last.SalaryTo == nil {
			return nil, vacancies
		}
		cursor = &SalaryCursor{
			SalaryFrom: *last.SalaryFrom,
			SalaryTo:   *last.SalaryTo,
			ID:         last.ID,
		}
	}

	c, _ := cursor.(*CursorT)
	return c, vacancies
}

func buildResponse[CursorT any](
	vacancies []views.PublishedVacSummary,
	nextCursor *CursorT,
	order Order,
) (*Response, error) {
	nextCursorEncoded, err := encoding.EncodeCursor[CursorT, Order](order, nextCursor)
	if err != nil {
		return nil, err
	}

	response := Response{
		Vacancies:  vacancies,
		NextCursor: nextCursorEncoded,
	}

	return &response, nil
}
