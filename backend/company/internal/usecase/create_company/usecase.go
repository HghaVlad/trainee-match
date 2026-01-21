package create_company

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
)

type Usecase struct {
	repo CompanyRepo
}

func NewUsecase(repo CompanyRepo) *Usecase {
	return &Usecase{repo}
}

func (u *Usecase) Execute(ctx context.Context, request *Request) (*Response, error) {

	// TODO: do smth with owner id

	company := &entities.Company{
		ID:          uuid.New(),
		Name:        request.Name,
		Description: request.Description,
		Website:     request.Website,
		OwnerID:     uuid.New(),
	}

	err := u.repo.Create(ctx, company)
	if err != nil {
		return nil, err
	}

	resp := &Response{
		ID: company.ID,
	}

	return resp, nil
}
