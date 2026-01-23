package auth

import "errors"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var ErrorInvalidToken = errors.New("invalid access token")
