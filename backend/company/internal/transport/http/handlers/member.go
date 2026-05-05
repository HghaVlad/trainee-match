package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/dto"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/helpers"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/mappers"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/middleware"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/add"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/list"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/remove"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/update"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/projection/userhr"
)

type MemberHandler struct {
	add    *add.Usecase
	list   *list.Usecase
	update *update.Usecase
	delete *remove.Usecase
}

func NewMemberHandler(
	add *add.Usecase,
	list *list.Usecase,
	update *update.Usecase,
	del *remove.Usecase,
) *MemberHandler {
	return &MemberHandler{
		add:    add,
		list:   list,
		update: update,
		delete: del,
	}
}

// Add godoc
// @Summary Add member to company
// @Description Adds member to company by username. Requires admin role in company
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
		expected := h.handleErr(ctx, w, err)
		if !expected {
			handleUnexpectedErr(ctx, w, err, "failed to add company member",
				"username", dtoReq.Username)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// List godoc
// @Summary List members of company.
// @Description List members of company with usernames, emails, sorted by username; requires being a member of the company. Standard limit / offset pagination, with hasMore.
// @Tags member
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param limit query int false "Items per page" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} dto.CompanyMemberListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /companies/{id}/members [get]
func (h *MemberHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	iden := middleware.IdentityFromContext(ctx)

	req := list.Request{
		CompanyID: middleware.UUIDFromContext(ctx, "company-id"),
		Limit:     helpers.ParseLimit(r, "limit", 20),
		Offset:    helpers.ParseLimit(r, "offset", 0),
	}

	res, err := h.list.Execute(ctx, req, iden)

	if err != nil {
		expected := h.handleErr(ctx, w, err)
		if !expected {
			handleUnexpectedErr(ctx, w, err, "failed to list company members",
				"company_id", req.CompanyID)
		}
		return
	}

	resp := mappers.CompanyMemberListRespToDto(*res)
	helpers.RespondJSON(ctx, w, http.StatusOK, resp)
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
		expected := h.handleErr(ctx, w, err)
		if !expected {
			handleUnexpectedErr(ctx, w, err, "failed to update company member",
				"member_id", userID)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Remove godoc
// @Summary Remove company member
// @Description removes company member. Requires admin role in company. Admin can't remove themselves if they are the only admin left.
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
// @Router /companies/{id}/members/{user-id} [delete]
func (h *MemberHandler) Remove(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	iden := middleware.IdentityFromContext(ctx)
	companyID := middleware.UUIDFromContext(ctx, "company-id")
	userID := middleware.UUIDFromContext(ctx, "user-id")

	err := h.delete.Execute(ctx, companyID, userID, iden)
	if err != nil {
		expected := h.handleErr(ctx, w, err)
		if !expected {
			handleUnexpectedErr(ctx, w, err, "failed to delete company member",
				"member_id", userID)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MemberHandler) handleErr(ctx context.Context, w http.ResponseWriter, err error) bool {
	switch {
	case errors.Is(err, company.ErrCompanyNotFound),
		errors.Is(err, member.ErrCompanyMemberNotFound),
		errors.Is(err, userhr.ErrNotFound):
		helpers.RespondError(ctx, w, http.StatusNotFound, err)
		return true

	case errors.Is(err, member.ErrCompanyMemberAlreadyExists):
		helpers.RespondError(ctx, w, http.StatusConflict, err)
		return true

	case errors.Is(err, member.ErrInvalidUserID),
		errors.Is(err, member.ErrInvalidCompanyMemberRole),
		errors.Is(err, member.ErrCantRemoveYourself),
		errors.Is(err, common.ErrLimitTooLarge),
		errors.Is(err, userhr.ErrUsernameEmpty):
		helpers.RespondError(ctx, w, http.StatusBadRequest, err)
		return true

	case errors.Is(err, identity.ErrHrRoleRequired),
		errors.Is(err, member.ErrCompanyMemberRequired),
		errors.Is(err, member.ErrInsufficientRoleInCompany):
		helpers.RespondError(ctx, w, http.StatusForbidden, err)
		return true

	case errors.Is(err, context.DeadlineExceeded):
		helpers.RespondErrorMsg(ctx, w, http.StatusGatewayTimeout, "timeout: operation took too long")
		return true

	default:
		return false
	}
}
