package list

import (
	"context"
	"time"

	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/services/encoding"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
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

	// limit + 1 strat
	vacancies, err := uc.repo.ListPublished(ctx, req.Requirements, req.Order, cursor, req.Limit+1)
	if err != nil {
		return nil, err
	}

	nextCursor, vacancies := getNextCursor[CursorT](vacancies, req.Limit)

	resp, err := buildResponse[CursorT](vacancies, nextCursor, req.Order)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func getNextCursor[CursorT any](vacancies []VacancySummary, limit int) (*CursorT, []VacancySummary) {
	if len(vacancies) <= limit {
		return nil, vacancies
	}

	vacancies = vacancies[:len(vacancies)-1]
	last := vacancies[len(vacancies)-1]

	var zero CursorT

	switch any(zero).(type) {

	case PublishedAtCursor:
		cursor := PublishedAtCursor{
			PublishedAt: last.PublishedAt,
			Id:          last.ID,
		}
		return any(&cursor).(*CursorT), vacancies

	case SalaryCursor:
		if last.SalaryFrom == nil || last.SalaryTo == nil {
			return nil, vacancies
		}
		cursor := SalaryCursor{
			SalaryFrom: *last.SalaryFrom,
			SalaryTo:   *last.SalaryTo,
			Id:         last.ID,
		}
		return any(&cursor).(*CursorT), vacancies
	}

	return nil, nil
}

func buildResponse[CursorT any](vacancies []VacancySummary, nextCursor *CursorT, order Order) (*Response, error) {
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
