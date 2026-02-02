package list_vacancy

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/services/encoding"
)

type Usecase struct {
	repo      VacancyRepo
	respCache ResponseCacheRepo
}

func NewUsecase(repo VacancyRepo, cache ResponseCacheRepo) *Usecase {
	return &Usecase{repo: repo, respCache: cache}
}

func (uc *Usecase) Execute(ctx context.Context, req *Request) (*Response, error) {
	resp := uc.getFromCache(ctx, req)
	if resp != nil {
		return resp, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	switch req.Order {
	case OrderPublishedAtDesc:
		vacancies, nextCursor, err := uc.listByPublishedAt(ctx, req)
		if err != nil {
			return nil, err
		}
		resp, err = buildResponse[PublishedAtCursor](vacancies, nextCursor, OrderPublishedAtDesc)
		if err != nil {
			return nil, err
		}

	default:
		return nil, domain_errors.ErrUnsupportedListOrder
	}

	uc.saveToCache(ctx, req, resp)
	return resp, nil
}

func (uc *Usecase) getFromCache(ctx context.Context, req *Request) *Response {
	respCacheKey := strings.Join([]string{
		string(req.Order), req.Cursor, strconv.Itoa(req.Limit),
	}, "-")

	return uc.respCache.Get(ctx, respCacheKey)
}

func (uc *Usecase) saveToCache(ctx context.Context, req *Request, resp *Response) {
	respCacheKey := strings.Join([]string{
		string(req.Order), req.Cursor, strconv.Itoa(req.Limit),
	}, "-")

	// Adding to cache with short ttl because it won't be updated/deleted by service
	uc.respCache.Put(ctx, respCacheKey, resp, time.Second*20)
}

func (uc *Usecase) listByPublishedAt(ctx context.Context, req *Request) ([]VacancySummary, *PublishedAtCursor, error) {
	cursor, curErr := encoding.DecodeCursor[PublishedAtCursor, Order](req.Cursor, req.Order)
	if curErr != nil {
		return nil, nil, curErr
	}

	return uc.repo.ListByPublishedAt(ctx, cursor, req.Limit)
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
