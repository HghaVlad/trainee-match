package services

import (
	"context"

	"github.com/Nerzal/gocloak/v13"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/domain"
)

type AuthRepo interface {
	CreateUser(ctx context.Context, user domain.User, password string) (string, error)
	Login(ctx context.Context, username, password string) (*gocloak.JWT, error)
	Logout(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error)
	GetUserInfo(ctx context.Context, token string) (*domain.User, error)
	GetUserRole(ctx context.Context, token string, userId string) (string, error)
}

type Auth struct {
	repo AuthRepo
}

func NewAuth(repo AuthRepo) *Auth {
	return &Auth{repo: repo}
}

func (a *Auth) Register(ctx context.Context, user domain.User, password string) (string, error) {
	return a.repo.CreateUser(ctx, user, password)
}

func (a *Auth) Login(ctx context.Context, username, password string) (*gocloak.JWT, error) {
	return a.repo.Login(ctx, username, password)
}

func (a *Auth) Logout(ctx context.Context, token string) error {
	return a.repo.Logout(ctx, token)
}

func (a *Auth) RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error) {
	return a.repo.RefreshToken(ctx, refreshToken)
}

func (a *Auth) GetUserMe(ctx context.Context, token string) (*domain.User, error) {
	user, err := a.repo.GetUserInfo(ctx, token)
	if err != nil {
		return nil, err
	}
	role, err := a.repo.GetUserRole(ctx, token, user.Id)
	if err != nil {
		return nil, err
	}
	user.Role = role
	return user, nil
}
