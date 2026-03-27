package handlers

import (
	"errors"
	"net/http"

	gmiddleware "github.com/M0s1ck/g-store/src/pkg/http/middleware"
	"github.com/M0s1ck/g-store/src/pkg/http/responds"

	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/helpers"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/mapper"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/middleware"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/add"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/delete"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/update"
)

type MemberHandler struct {
	add    *add.Usecase
	update *update.Usecase
	delete *delete.Usecase
}

func NewMemberHandler(
	add *add.Usecase,
	update *update.Usecase,
	del *delete.Usecase,
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
// @Failure 400 {object} responds.ErrorResponse
// @Failure 401 {object} responds.ErrorResponse
// @Failure 403 {object} responds.ErrorResponse
// @Failure 404 {object} responds.ErrorResponse
// @Failure 409 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /companies/{id}/members [post]
func (h *MemberHandler) Add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	iden := middleware.IdentityFromContext(ctx)

	id, err := gmiddleware.UUIDFromContext(ctx)
	if err != nil {
		responds.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	dtoReq, err := gmiddleware.BodyFromContext[dto.CompanyAddHrRequest](ctx)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	req := mapper.CompanyAddHrReqToUC(id, dtoReq)

	err = h.add.Execute(ctx, req, iden)
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
// @Failure 400 {object} responds.ErrorResponse
// @Failure 401 {object} responds.ErrorResponse
// @Failure 403 {object} responds.ErrorResponse
// @Failure 404 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /companies/{id}/members/{user-id} [patch]
func (h *MemberHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	iden := middleware.IdentityFromContext(ctx)

	companyID, err := gmiddleware.UUIDFromContext(ctx)
	if err != nil {
		responds.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	userID, ok := helpers.ParseUuidFromPathOr400(r, w, "user-id")
	if !ok {
		return
	}

	dtoReq, err := gmiddleware.BodyFromContext[dto.CompanyUpdateMemberRequest](ctx)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	req := mapper.CompanyUpdateMemberReqToUC(companyID, userID, dtoReq)

	err = h.update.Execute(ctx, req, iden)
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
// @Failure 400 {object} responds.ErrorResponse
// @Failure 401 {object} responds.ErrorResponse
// @Failure 403 {object} responds.ErrorResponse
// @Failure 404 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /companies/{id}/members/{user-id} [delete]
func (h *MemberHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	iden := middleware.IdentityFromContext(ctx)

	companyID, err := gmiddleware.UUIDFromContext(ctx)
	if err != nil {
		responds.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	userID, ok := helpers.ParseUuidFromPathOr400(r, w, "user-id")
	if !ok {
		return
	}

	err = h.delete.Execute(ctx, companyID, userID, iden)
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
		responds.RespondError(w, http.StatusNotFound, err)

	case errors.Is(err, member.ErrCompanyMemberAlreadyExists):
		responds.RespondError(w, http.StatusConflict, err)

	case errors.Is(err, member.ErrInvalidUserID),
		errors.Is(err, member.ErrInvalidCompanyMemberRole):
		responds.RespondError(w, http.StatusBadRequest, err)

	case errors.Is(err, identity.ErrHrRoleRequired),
		errors.Is(err, member.ErrCompanyMemberRequired),
		errors.Is(err, member.ErrInsufficientRoleInCompany):
		responds.RespondError(w, http.StatusForbidden, err)

	default:
		responds.RespondError(w, http.StatusInternalServerError, err)
	}
}
