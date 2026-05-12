package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Nerzal/gocloak/v13"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/delivery/http/helpers"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/domain"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/usecase/get_user_me"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/usecase/login"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/usecase/logout"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/usecase/refresh_token"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/usecase/register"
)

type RegisterUseCase interface {
	Execute(ctx context.Context, req *register.Request) (uuid.UUID, error)
}

type LoginUseCase interface {
	Execute(ctx context.Context, req *login.Request) (*gocloak.JWT, error)
}

type LogoutUseCase interface {
	Execute(ctx context.Context, req *logout.Request) error
}

type RefreshTokenUseCase interface {
	Execute(ctx context.Context, req *refresh_token.Request) (*gocloak.JWT, error)
}

type GetUserMeUseCase interface {
	Execute(ctx context.Context, req *get_user_me.Request) (*domain.User, error)
}

type Auth struct {
	registerUC          RegisterUseCase
	loginUC             LoginUseCase
	logoutUC            LogoutUseCase
	refreshTokenUC      RefreshTokenUseCase
	getUserMeUC         GetUserMeUseCase
	validate            *validator.Validate
	AccessTokenExpires  int
	RefreshTokenExpires int
}

func NewAuthHandler(
	registerUC RegisterUseCase,
	loginUC LoginUseCase,
	logoutUC LogoutUseCase,
	refreshTokenUC RefreshTokenUseCase,
	getUserMeUC GetUserMeUseCase,
	accessTokenExpires,
	refreshTokenExpires int,
) *Auth {
	return &Auth{
		registerUC:          registerUC,
		loginUC:             loginUC,
		logoutUC:            logoutUC,
		refreshTokenUC:      refreshTokenUC,
		getUserMeUC:         getUserMeUC,
		validate:            validator.New(),
		AccessTokenExpires:  accessTokenExpires,
		RefreshTokenExpires: refreshTokenExpires,
	}
}

// Register godoc
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.RegisterUserRequest true "User registration data"
// @Success 200 {object} domain.User
// @Failure 400 {object} dto.ErrorResponse "invalid request"
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/register [post]
func (h *Auth) Register(w http.ResponseWriter, r *http.Request) {
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

	id, err := h.registerUC.Execute(r.Context(), &register.Request{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Username:  request.Username,
		Password:  request.Password,
		Role:      request.Role,
	})
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	user.ID = id.String()

	helpers.RespondJSON(w, http.StatusOK, user)
}

// Login godoc
// @Summary Login a user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.LoginRequest true "User login data"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse "invalid request"
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *Auth) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request dto.LoginRequest
	err := decoder.Decode(&request)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}

	token, err := h.loginUC.Execute(r.Context(), &login.Request{
		Username: request.Username,
		Password: request.Password,
	})

	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	helpers.SetTokenPairToCookies(w, token.AccessToken, token.RefreshToken, h.AccessTokenExpires, h.RefreshTokenExpires)

	helpers.RespondJSON(w, http.StatusOK, dto.MessageResponse{Message: "OK"})
}

// RefreshToken godoc
// @Summary Refresh JWT token
// @Tags auth
// @Produce json
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/refresh [post]
func (h *Auth) RefreshToken(w http.ResponseWriter, r *http.Request) {
	token := helpers.GetRefreshTokenFromCookies(r)
	if token == "" {
		helpers.RespondError(w, http.StatusBadRequest, "missing refresh token")
		return
	}

	newToken, err := h.refreshTokenUC.Execute(r.Context(), &refresh_token.Request{
		RefreshToken: token,
	})
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	helpers.SetTokenPairToCookies(
		w,
		newToken.AccessToken,
		newToken.RefreshToken,
		h.AccessTokenExpires,
		h.RefreshTokenExpires,
	)

	helpers.RespondJSON(w, http.StatusOK, dto.MessageResponse{Message: "OK"})
}

// Logout godoc
// @Summary Logout a user
// @Tags auth
// @Produce json
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/logout [post]
func (h *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	refreshToken := helpers.GetRefreshTokenFromCookies(r)
	helpers.SetTokenPairToCookies(w, "", "", 0, 0)
	if refreshToken == "" {
		helpers.RespondJSON(w, http.StatusOK, dto.MessageResponse{Message: "Successfully log out"})
		return
	}
	_ = h.logoutUC.Execute(r.Context(), &logout.Request{Token: refreshToken})
	helpers.RespondJSON(w, http.StatusOK, dto.MessageResponse{Message: "Successfully log out"})
}

func (h *Auth) GetMe(w http.ResponseWriter, r *http.Request) {
	token := helpers.GetAccessTokenFromCookies(r)
	if token == "" {
		helpers.RespondError(w, http.StatusBadRequest, "missing access token")
		return
	}
	user, err := h.getUserMeUC.Execute(r.Context(), &get_user_me.Request{Token: token})
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	}
	helpers.RespondJSON(w, http.StatusOK, response)
}
