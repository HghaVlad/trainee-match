package services

import (
	"context"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/domain"
	"github.com/Nerzal/gocloak/v13"
)

type AuthRepo interface {
	CreateUser(ctx context.Context, user domain.User, password string) (string, error)
	Login(ctx context.Context, username, password string) (*gocloak.JWT, error)
	Logout(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error)
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
