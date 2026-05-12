package getuser

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/domain"
)

//go:generate mockery
type AuthRepo interface {
	GetUserInfo(ctx context.Context, token string) (*domain.User, error)
	GetUserRole(ctx context.Context, token string, userID string) (string, error)
}

type UseCase struct {
	repo AuthRepo
}

func New(repo AuthRepo) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Execute(ctx context.Context, req *Request) (*domain.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if req == nil {
		return nil, context.Canceled
	}
	user, err := uc.repo.GetUserInfo(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	role, err := uc.repo.GetUserRole(ctx, req.Token, user.ID)
	if err != nil {
		return nil, err
	}
	user.Role = role
	return user, nil
}
