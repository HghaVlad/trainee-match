package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/auth"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/delivery/http/helpers"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/domain"
	"github.com/Nerzal/gocloak/v13"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type AuthClient interface {
	CreateUser(ctx context.Context, user domain.User, password string) (string, error)
	Login(ctx context.Context, request auth.LoginRequest) (*gocloak.JWT, error)
	Logout(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error)
}

type AuthHandler struct {
	authClient AuthClient
	validate   *validator.Validate
}

func NewAuthHandler(authClient AuthClient) *AuthHandler {
	return &AuthHandler{
		authClient: authClient,
		validate:   validator.New(),
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request dto.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}
	if err := h.validate.Struct(request); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}

	user := domain.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Username:  request.Username,
		Role:      request.Role,
	}

	id, err := h.authClient.CreateUser(r.Context(), user, request.Password)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	user.Id = id

	helpers.RespondJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request dto.LoginRequest
	err := decoder.Decode(&request)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}

	token, err := h.authClient.Login(r.Context(), auth.LoginRequest{
		Username: request.Username,
		Password: request.Password,
	})

	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	helpers.RespondJSON(w, http.StatusOK, dto.TokenPair{AccessToken: token.AccessToken, RefreshToken: token.RefreshToken})
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request dto.RefreshTokenRequest
	err := decoder.Decode(&request)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}

	newToken, err := h.authClient.RefreshToken(r.Context(), request.RefreshToken)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	helpers.RespondJSON(w, http.StatusOK, dto.TokenPair{AccessToken: newToken.AccessToken, RefreshToken: newToken.RefreshToken})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		helpers.RespondError(w, http.StatusBadRequest, "missing authorization header")
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	err := h.authClient.Logout(r.Context(), token)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	helpers.RespondJSON(w, http.StatusOK, dto.MessageResponse{Message: "Successfully log out"})
}
