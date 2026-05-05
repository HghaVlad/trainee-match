package userhr

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Projection struct {
	UserID    uuid.UUID
	Username  string
	Email     string
	CreatedAt time.Time
}

var (
	ErrUserIDNil     = errors.New("user id cannot be nil")
	ErrUsernameEmpty = errors.New("user name cannot be empty")
	ErrEmailEmpty    = errors.New("email cannot be empty")
)

func NewHrProjection(userID uuid.UUID, username string, email string, createdAt time.Time) (*Projection, error) {
	if userID == uuid.Nil {
		return nil, ErrUserIDNil
	}

	if username == "" {
		return nil, ErrUsernameEmpty
	}

	if email == "" {
		return nil, ErrEmailEmpty
	}

	return &Projection{
		UserID:    userID,
		Username:  username,
		Email:     email,
		CreatedAt: createdAt,
	}, nil
}
