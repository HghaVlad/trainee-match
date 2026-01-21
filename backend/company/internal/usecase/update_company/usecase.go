package update_company

import "context"

type Usecase struct {
	repo CompanyRepo
}

func NewUsecase(repo CompanyRepo) *Usecase {
	return &Usecase{repo}
}

func (u *Usecase) Execute(ctx context.Context, req *Request) error {
	return u.repo.Update(ctx, req)
}
