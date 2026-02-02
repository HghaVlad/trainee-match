package list_companies

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
)

type Usecase struct {
	repo          Repo
	responseCache ResponseCacheRepo
}

func NewUsecase(repo Repo, responseCache ResponseCacheRepo) *Usecase {
	return &Usecase{
		repo:          repo,
		responseCache: responseCache,
	}
}

func (u *Usecase) Execute(ctx context.Context, req *Request) (*Response, error) {

	respCacheKey := strings.Join([]string{
		string(req.Order), req.Cursor, strconv.Itoa(req.Limit),
	}, "-")

	resp := u.responseCache.Get(ctx, respCacheKey)
	if resp != nil {
		return resp, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var err error

	switch req.Order {
	case OrderVacanciesDesc:
		resp, err = u.listByVacanciesCnt(ctx, req)

	case OrderCreatedAtDesc:
		resp, err = u.listByCreatedAt(ctx, req)

	case OrderNameAsc:
		resp, err = u.listByName(ctx, req)

	default:
		return nil, domain_errors.ErrUnsupportedListOrder
	}

	if err != nil {
		return nil, err
	}

	// Adding to cache with short ttl because it won't be updated/deleted by service
	u.responseCache.Put(ctx, respCacheKey, resp, time.Second*20)
	return resp, nil
}

func (u *Usecase) listByCreatedAt(ctx context.Context, req *Request) (*Response, error) {
	cursor, curErr := decodeCursor[CreatedAtCursor](req.Cursor, req.Order)
	if curErr != nil {
		return nil, curErr
	}

	companies, nextCursor, err := u.repo.ListByCreatedAtDesc(ctx, cursor, req.Limit)
	if err != nil {
		return nil, err
	}

	nextCursorEncoded, err := encodeCursor[CreatedAtCursor](OrderCreatedAtDesc, nextCursor)
	if err != nil {
		return nil, err
	}

	response := Response{
		Companies:  companies,
		NextCursor: nextCursorEncoded,
	}

	return &response, nil
}

func (u *Usecase) listByVacanciesCnt(ctx context.Context, req *Request) (*Response, error) {
	cursor, curErr := decodeCursor[VacanciesCntCursor](req.Cursor, req.Order)
	if curErr != nil {
		return nil, curErr
	}

	companies, nextCursor, err := u.repo.ListByVacanciesCnt(ctx, cursor, req.Limit)
	if err != nil {
		return nil, err
	}

	nextCursorEncoded, err := encodeCursor[VacanciesCntCursor](OrderVacanciesDesc, nextCursor)
	if err != nil {
		return nil, err
	}

	response := Response{
		Companies:  companies,
		NextCursor: nextCursorEncoded,
	}

	return &response, nil
}

func (u *Usecase) listByName(ctx context.Context, req *Request) (*Response, error) {
	cursor, curErr := decodeCursor[NameCursor](req.Cursor, req.Order)
	if curErr != nil {
		return nil, curErr
	}

	companies, nextCursor, err := u.repo.ListByName(ctx, cursor, req.Limit)
	if err != nil {
		return nil, err
	}

	nextCursorEncoded, err := encodeCursor[NameCursor](OrderNameAsc, nextCursor)
	if err != nil {
		return nil, err
	}

	response := Response{
		Companies:  companies,
		NextCursor: nextCursorEncoded,
	}

	return &response, nil
}
