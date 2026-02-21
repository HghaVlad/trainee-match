package domain

import (
	"errors"
	"github.com/google/uuid"
	"regexp"
	"time"
)

var (
	ErrCandidateNotFound     = errors.New("candidate not found")
	ErrBirthdayInFuture      = errors.New("birthday cannot be in the future")
	ErrInvalidPhoneFormat    = errors.New("invalid phone number format")
	ErrInvalidTelegramFormat = errors.New("telegram username must start with @ and contain only alphanumeric characters and underscores, 3-32 characters long")
	ErrInvalidCityFormat     = errors.New("city is required")
	ErrTelegramAlreadyExists = errors.New("telegram username already exists")
	ErrPhoneAlreadyExists    = errors.New("phone number already exists")
)

type Candidate struct {
	ID       uuid.UUID
	UserId   uuid.UUID
	Phone    string
	Telegram string
	City     string
	Birthday time.Time
}

func (c Candidate) Validate() error {
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	if c.Phone == "" || !phoneRegex.MatchString(c.Phone) {
		return ErrInvalidPhoneFormat
	}

	telegramRegex := regexp.MustCompile(`^@[\w]{3,32}$`)
	if c.Telegram == "" || !telegramRegex.MatchString(c.Telegram) {
		return ErrInvalidTelegramFormat
	}

	if c.City == "" {
		return ErrInvalidCityFormat
	}

	if c.Birthday.After(time.Now()) {
		return ErrBirthdayInFuture
	}

	return nil
}
