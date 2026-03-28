package list

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/utils/encoding"
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
	if err := req.Validate(); err != nil {
		return nil, err
	}

	respCacheKey := strings.Join([]string{
		string(req.Order), req.EncodedCursor, strconv.Itoa(req.Limit),
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
		resp, err = list[VacanciesCntCursor](ctx, u, req)
	case OrderCreatedAtDesc:
		resp, err = list[CreatedAtCursor](ctx, u, req)
	case OrderNameAsc:
		resp, err = list[NameCursor](ctx, u, req)

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

func list[CursorT any](ctx context.Context, uc *Usecase, req *Request) (*Response, error) {
	cursor, curErr := encoding.DecodeCursor[CursorT, Order](req.EncodedCursor, req.Order)
	if curErr != nil {
		return nil, curErr
	}

	// limit + 1 strat
	companies, err := uc.repo.ListSummaries(ctx, req.Order, cursor, req.Limit+1)
	if err != nil {
		return nil, err
	}

	nextCursor, companies := getNextCursor[CursorT](companies, req.Limit)

	resp, err := buildResponse[CursorT](companies, nextCursor, req.Order)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func getNextCursor[CursorT any](companies []CompanySummary, limit int) (*CursorT, []CompanySummary) {
	if len(companies) <= limit {
		return nil, companies
	}

	companies = companies[:len(companies)-1]
	last := companies[len(companies)-1]

	var zero CursorT
	var cursor any

	switch any(zero).(type) {
	case VacanciesCntCursor:
		cursor = &VacanciesCntCursor{
			Count: last.OpenVacanciesCnt,
			Name:  last.Name,
		}

	case CreatedAtCursor:
		cursor = &CreatedAtCursor{
			CreatedAt: last.CreatedAt,
			Name:      last.Name,
		}

	case NameCursor:
		cursor = &NameCursor{
			Name: last.Name,
		}
	}

	c, _ := cursor.(*CursorT)
	return c, companies
}

func buildResponse[CursorT any](companies []CompanySummary, nextCursor *CursorT, order Order) (*Response, error) {
	nextCursorEncoded, err := encoding.EncodeCursor[CursorT, Order](order, nextCursor)
	if err != nil {
		return nil, err
	}

	response := Response{
		Companies:  companies,
		NextCursor: nextCursorEncoded,
	}

	return &response, nil
}
