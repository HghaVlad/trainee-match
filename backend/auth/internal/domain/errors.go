package domain

import "errors"

var (
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidName        = errors.New("invalid name")
	ErrEmailAlreadyExists = errors.New("email already exists")
)
