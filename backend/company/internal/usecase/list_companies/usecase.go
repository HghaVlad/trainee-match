package list_companies

import (
	"context"
	"errors"
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

	// TODO: перевыбрасывать везде свои ошибки и хендлить их
	// TODO: if wait 10 sec here, writer won't answer, check and solve this

	switch req.Order {
	case OrderVacanciesDesc:
		return nil, nil

	case OrderCreatedAtDesc:
		return u.listByCreatedAt(ctx, req)

	case OrderNameAsc:
		return nil, nil

	default:
		return nil, errors.New("unsupported list order")
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
