package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/delivery/http/helpers"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/usecase/admin/newadmin"
)

type AddAdminRoleUsecase interface {
	Execute(ctx context.Context, req newadmin.Request) error
}
type Admin struct {
	addAdminRoleUsecase AddAdminRoleUsecase
}

func NewAdmin(addAdminRoleUsecase AddAdminRoleUsecase) *Admin {
	return &Admin{addAdminRoleUsecase: addAdminRoleUsecase}
}

// NewAdmin godoc
// @Summary Add admin role to new user
// @Tags admin
// @Accept json
// @Produce json
// @Param input body dto.AddAdminRoleRequest true "Make new admin"
// @Success 200 {object} dto.MessageResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/new [post]
func (admin *Admin) NewAdmin(w http.ResponseWriter, r *http.Request) {
	token := helpers.GetAccessTokenFromCookies(r)
	if token == "" {
		helpers.RespondError(w, http.StatusUnauthorized, "missing access token")
		return
	}

	var req dto.AddAdminRoleRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}

	err := admin.addAdminRoleUsecase.Execute(r.Context(), newadmin.Request{UserID: req.UserID, AccessToken: token})
	if err != nil {
		helpers.RespondSmartError(w, err)
		return
	}
	helpers.RespondJSON(w, http.StatusOK, dto.MessageResponse{Message: "OK"})
}
