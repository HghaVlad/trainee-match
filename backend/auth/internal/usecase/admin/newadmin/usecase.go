package newadmin

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/domain"
)

type Repo interface {
	GetUserInfo(ctx context.Context, token string) (*domain.User, error)
	GetUserRole(ctx context.Context, _ string, userID string) (string, error)
	AddAdminRole(ctx context.Context, userID string) error
}

type UseCase struct {
	Repo Repo
}

func NewUseCase(repo Repo) *UseCase {
	return &UseCase{Repo: repo}
}

func (uc *UseCase) Execute(ctx context.Context, req Request) error {
	curUser, err := uc.Repo.GetUserInfo(ctx, req.AccessToken)
	if err != nil {
		return err
	}
	userRole, err := uc.Repo.GetUserRole(ctx, "", curUser.ID)
	if err != nil {
		return domain.ErrForbidden
	}
	if domain.UserRole(userRole) != domain.UserAdminRole {
		return domain.ErrForbidden
	}

	return uc.Repo.AddAdminRole(ctx, req.UserID)
}
