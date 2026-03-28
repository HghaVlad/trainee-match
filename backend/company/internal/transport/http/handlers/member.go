package handlers

import (
	"errors"
	"net/http"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/dto"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/helpers"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/mappers"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/middleware"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/add"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/remove"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/update"
)

type MemberHandler struct {
	add    *add.Usecase
	update *update.Usecase
	delete *remove.Usecase
}

func NewMemberHandler(
	add *add.Usecase,
	update *update.Usecase,
	del *remove.Usecase,
) *MemberHandler {
	return &MemberHandler{
		add:    add,
		update: update,
		delete: del,
	}
}

// Add godoc
// @Summary Add member to company
// @Description Adds member to company. Requires admin role in company
// @Tags member
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param company_add_hr_request body dto.CompanyAddHrRequest true "Request to add member"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /companies/{id}/members [post]
func (h *MemberHandler) Add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	iden := middleware.IdentityFromContext(ctx)
	compID := middleware.UUIDFromContext(ctx, "company-id")
	dtoReq := middleware.BodyFromContext[dto.CompanyAddHrRequest](ctx)

	req := mappers.CompanyAddHrReqToUC(compID, dtoReq)

	err := h.add.Execute(ctx, req, iden)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Update godoc
// @Summary Update company member
// @Description Updates only company member role. Requires admin role in company
// @Tags member
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param user-id path string true "User ID"
// @Param company_update_member_request body dto.CompanyUpdateMemberRequest true "Request to update member"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /companies/{id}/members/{user-id} [patch]
func (h *MemberHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	iden := middleware.IdentityFromContext(ctx)
	companyID := middleware.UUIDFromContext(ctx, "company-id")
	userID := middleware.UUIDFromContext(ctx, "user-id")
	dtoReq := middleware.BodyFromContext[dto.CompanyUpdateMemberRequest](ctx)

	req := mappers.CompanyUpdateMemberReqToUC(companyID, userID, dtoReq)

	err := h.update.Execute(ctx, req, iden)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Delete godoc
// @Summary Delete company member
// @Description Deletes company member. Requires admin role in company
// @Tags member
// @Produce json
// @Param id path string true "Company ID"
// @Param user-id path string true "User ID"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /companies/{id}/members/{user-id} [remove]
func (h *MemberHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	iden := middleware.IdentityFromContext(ctx)
	companyID := middleware.UUIDFromContext(ctx, "company-id")
	userID := middleware.UUIDFromContext(ctx, "user-id")

	err := h.delete.Execute(ctx, companyID, userID, iden)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MemberHandler) handleErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, company.ErrCompanyNotFound),
		errors.Is(err, member.ErrCompanyMemberNotFound):
		helpers.RespondError(w, http.StatusNotFound, err)

	case errors.Is(err, member.ErrCompanyMemberAlreadyExists):
		helpers.RespondError(w, http.StatusConflict, err)

	case errors.Is(err, member.ErrInvalidUserID),
		errors.Is(err, member.ErrInvalidCompanyMemberRole):
		helpers.RespondError(w, http.StatusBadRequest, err)

	case errors.Is(err, identity.ErrHrRoleRequired),
		errors.Is(err, member.ErrCompanyMemberRequired),
		errors.Is(err, member.ErrInsufficientRoleInCompany):
		helpers.RespondError(w, http.StatusForbidden, err)

	default:
		helpers.RespondError(w, http.StatusInternalServerError, err)
	}
}
