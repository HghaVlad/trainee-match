package register

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/domain"
)

type AuthRepo interface {
	CreateUser(ctx context.Context, user domain.User, password string) (string, error)
}

type OutboxWriter interface {
	WriteUserCreated(ctx context.Context, ev domain.UserCreatedEvent) error
}

type UseCase struct {
	repo         AuthRepo
	outboxWriter OutboxWriter
	validate     *validator.Validate
}

func New(repo AuthRepo, outboxWriter OutboxWriter) *UseCase {
	return &UseCase{
		repo:         repo,
		outboxWriter: outboxWriter,
		validate:     validator.New(),
	}
}

func (uc *UseCase) Execute(ctx context.Context, req *Request) (uuid.UUID, error) {
	if err := req.Validate(uc.validate); err != nil {
		return uuid.Nil, err
	}

	user := domain.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Username:  req.Username,
		Role:      req.Role,
	}

	id, err := uc.repo.CreateUser(ctx, user, req.Password)
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse user id: %w", err)
	}

	if uc.outboxWriter != nil {
		ev := domain.UserCreatedEvent{
			EventID:    uuid.New(),
			UserID:     userID,
			Username:   user.Username,
			Email:      user.Email,
			Role:       user.Role,
			OccurredAt: time.Now().UTC(),
		}
		if err := uc.outboxWriter.WriteUserCreated(ctx, ev); err != nil {
			return uuid.Nil, err
		}
	}

	return userID, nil
}
