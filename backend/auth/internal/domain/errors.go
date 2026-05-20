package domain

import "errors"

var (
	ErrInvalidEmail          = errors.New("invalid email")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrInvalidName           = errors.New("invalid name")
	ErrInvalidRole           = errors.New("invalid role")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUserNameAlreadyExists = errors.New("username already exists")
	ErrIncorrectPassword     = errors.New("incorrect password")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrForbidden             = errors.New("forbidden")
)
