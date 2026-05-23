package login

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
)

type AuthRepo interface {
	Login(ctx context.Context, username, password string) (*gocloak.JWT, error)
}

type UseCase struct {
	repo AuthRepo
}

func New(repo AuthRepo) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Execute(ctx context.Context, req *Request) (*gocloak.JWT, error) {
	return uc.repo.Login(ctx, req.Username, req.Password)
}
