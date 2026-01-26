package list_companies

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
)

type Usecase struct {
	repo Repo
}

func NewUsecase(repo Repo) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) Execute(ctx context.Context, req *Request) (*Response, error) {
	// TODO: if wait 10 sec here, writer won't answer, check and solve this

	switch req.Order {
	case OrderVacanciesDesc:
		return u.listByVacanciesCnt(ctx, req)

	case OrderCreatedAtDesc:
		return u.listByCreatedAt(ctx, req)

	case OrderNameAsc:
		return u.listByName(ctx, req)

	default:
		return nil, domain_errors.ErrUnsupportedListOrder
	}
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
