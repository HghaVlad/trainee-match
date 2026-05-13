package logout

import "context"

type AuthRepo interface {
	Logout(ctx context.Context, token string) error
}

type UseCase struct {
	repo AuthRepo
}

func New(repo AuthRepo) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Execute(ctx context.Context, req *Request) error {
	return uc.repo.Logout(ctx, req.Token)
}
