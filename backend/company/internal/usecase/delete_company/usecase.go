package delete_company

import (
	"context"

	"github.com/google/uuid"
)

type Usecase struct {
	repo CompanyRepo
}

func NewUsecase(repo CompanyRepo) *Usecase {
	return &Usecase{repo}
}

func (u *Usecase) Execute(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
