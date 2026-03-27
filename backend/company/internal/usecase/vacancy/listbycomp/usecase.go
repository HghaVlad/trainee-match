package listbycomp

import (
	"context"
	"time"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/services/encoding"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
)

type Usecase struct {
	vacRepo   VacancyRepo
	compRepo  CompanyRepo
	respCache ResponseCacheRepo
}

func NewUsecase(vacRepo VacancyRepo, compRepo CompanyRepo, cache ResponseCacheRepo) *Usecase {
	return &Usecase{vacRepo: vacRepo, compRepo: compRepo, respCache: cache}
}

func (uc *Usecase) Execute(ctx context.Context, req *Request) (*Response, error) {
	respCacheKey := req.toCacheKey()
	resp := uc.respCache.Get(ctx, respCacheKey)
	if resp != nil {
		return resp, nil
	}

	companyExists, exErr := uc.compRepo.Exists(ctx, req.CompID)
	if exErr != nil {
		return nil, exErr
	}

	if !companyExists {
		return nil, company.ErrCompanyNotFound
	}

	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	var err error

	switch req.Order {
	case OrderPublishedAtDesc:
		resp, err = uc.listByPublishedAt(ctx, req)

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

func (uc *Usecase) listByPublishedAt(ctx context.Context, req *Request) (*Response, error) {
	cursor, curErr := encoding.DecodeCursor[PublishedAtCursor, Order](req.Cursor, req.Order)
	if curErr != nil {
		return nil, curErr
	}

	vacancies, err := uc.vacRepo.ListByCompanyByPublishedAt(ctx, req.CompID, cursor, req.Limit)
	if err != nil {
		return nil, err
	}

	// next cursor if full page
	var nextCursor *PublishedAtCursor = nil
	if len(vacancies) == req.Limit {
		last := vacancies[len(vacancies)-1]
		nextCursor = &PublishedAtCursor{
			PublishedAt: last.PublishedAt,
			Id:          last.ID,
		}
	}

	resp, err := buildResponse[PublishedAtCursor](vacancies, nextCursor, OrderPublishedAtDesc)
	if err != nil {
		return nil, err
	}

	return resp, nil
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
