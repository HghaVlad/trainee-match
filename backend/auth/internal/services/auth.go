package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/google/uuid"

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

type OutboxWriter interface {
	WriteUserCreated(ctx context.Context, ev domain.UserCreatedEvent) error
}

type Auth struct {
	repo         AuthRepo
	outboxWriter OutboxWriter
}

func NewAuth(repo AuthRepo, outboxWriter OutboxWriter) *Auth {
	return &Auth{repo: repo, outboxWriter: outboxWriter}
}

func (a *Auth) Register(ctx context.Context, user domain.User, password string) (string, error) {
	id, err := a.repo.CreateUser(ctx, user, password)
	if err != nil {
		return "", err
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		return "", fmt.Errorf("parse user id: %w", err)
	}

	if a.outboxWriter != nil {
		ev := domain.UserCreatedEvent{
			EventID:    uuid.New(),
			UserID:     userID,
			Username:   user.Username,
			Email:      user.Email,
			Role:       user.Role,
			OccurredAt: time.Now().UTC(),
		}
		if err := a.outboxWriter.WriteUserCreated(ctx, ev); err != nil {
			return "", err
		}
	}

	return id, nil
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
