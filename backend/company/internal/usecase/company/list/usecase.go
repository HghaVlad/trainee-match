package list

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/services/encoding"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
)

// Usecase cursor pagination list company. Supports different orders.
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

// Execute cursor pagination list company. Supports different orders
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
		resp, err = u.ListByVacanciesCnt(ctx, req)

	case OrderCreatedAtDesc:
		resp, err = u.ListByCreatedAt(ctx, req)

	case OrderNameAsc:
		resp, err = u.ListByName(ctx, req)

	default:
		return nil, common.ErrUnsupportedListOrder
	}

	if err != nil {
		return nil, err
	}

	// Adding to cache with short ttl because it won't be updated/deleted by service
	u.responseCache.Put(ctx, respCacheKey, resp, time.Second*20)
	return resp, nil
}

func (u *Usecase) ListByCreatedAt(ctx context.Context, req *Request) (*Response, error) {
	cursor, curErr := encoding.DecodeCursor[CreatedAtCursor, Order](req.Cursor, req.Order)
	if curErr != nil {
		return nil, curErr
	}

	companies, nextCursor, err := u.repo.ListByCreatedAtDesc(ctx, cursor, req.Limit)
	if err != nil {
		return nil, err
	}

	nextCursorEncoded, err := encoding.EncodeCursor[CreatedAtCursor, Order](OrderCreatedAtDesc, nextCursor)
	if err != nil {
		return nil, err
	}

	response := Response{
		Companies:  companies,
		NextCursor: nextCursorEncoded,
	}

	return &response, nil
}

func (u *Usecase) ListByVacanciesCnt(ctx context.Context, req *Request) (*Response, error) {
	cursor, curErr := encoding.DecodeCursor[VacanciesCntCursor, Order](req.Cursor, req.Order)
	if curErr != nil {
		return nil, curErr
	}

	companies, nextCursor, err := u.repo.ListByVacanciesCnt(ctx, cursor, req.Limit)
	if err != nil {
		return nil, err
	}

	nextCursorEncoded, err := encoding.EncodeCursor[VacanciesCntCursor, Order](OrderVacanciesDesc, nextCursor)
	if err != nil {
		return nil, err
	}

	response := Response{
		Companies:  companies,
		NextCursor: nextCursorEncoded,
	}

	return &response, nil
}

func (u *Usecase) ListByName(ctx context.Context, req *Request) (*Response, error) {
	cursor, curErr := encoding.DecodeCursor[NameCursor, Order](req.Cursor, req.Order)
	if curErr != nil {
		return nil, curErr
	}

	companies, nextCursor, err := u.repo.ListByName(ctx, cursor, req.Limit)
	if err != nil {
		return nil, err
	}

	nextCursorEncoded, err := encoding.EncodeCursor[NameCursor, Order](OrderNameAsc, nextCursor)
	if err != nil {
		return nil, err
	}

	response := Response{
		Companies:  companies,
		NextCursor: nextCursorEncoded,
	}

	return &response, nil
}
