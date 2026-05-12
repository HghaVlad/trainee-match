package refresh_token

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
)

type AuthRepo interface {
	RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error)
}

type UseCase struct {
	repo AuthRepo
}

func New(repo AuthRepo) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Execute(ctx context.Context, req *Request) (*gocloak.JWT, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, context.Canceled
	}
	return uc.repo.RefreshToken(ctx, req.RefreshToken)
}
