package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	utilslog "github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/utils/logger"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/dto"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/helpers"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/mappers"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/middleware"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/create"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/get"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/list"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/remove"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/update"
)

type CompanyHandler struct {
	getByID *get.Usecase
	create  *create.Usecase
	list    *list.Usecase
	update  *update.Usecase
	delete  *remove.Usecase
}

func NewCompanyHandler(
	get *get.Usecase,
	create *create.Usecase,
	list *list.Usecase,
	upd *update.Usecase,
	del *remove.Usecase,
) *CompanyHandler {
	return &CompanyHandler{
		getByID: get,
		create:  create,
		list:    list,
		update:  upd,
		delete:  del,
	}
}

// GetByID godoc
// @Summary Get profile by id
// @Description Returns company profile by UUID
// @Tags company
// @Accept json
// @Produce json
// @Param id path string true "Company ID (UUID)"
// @Success 200 {object} dto.CompanyResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /companies/{id} [get]
func (h *CompanyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := middleware.UUIDFromContext(ctx, "id")

	comp, err := h.getByID.Execute(ctx, id)

	if err != nil {
		expected := h.handleErr(ctx, w, err)
		if !expected {
			handleUnexpectedErr(ctx, w, err, "failed to get company", "id", id)
		}
		return
	}

	resp := mappers.GetCompRespToDto(comp)
	helpers.RespondJSON(ctx, w, http.StatusOK, resp)
}

// List godoc
// @Summary List company summaries
// @Description Uses cursor pagination, returns next cursor if there's more. Supports order by vacancies_desc, created_at_desc, name_asc
// @Tags company
// @Accept json
// @Produce json
// @Param order query string false "Order attribute" default(vacancies_desc)
// @Param cursor query string false "Items per page"
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.CompanyListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /companies [get]
func (h *CompanyHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limit := helpers.ParseLimit(r, "limit", 20)
	order := h.parseOrderQuery(r)
	cursor := r.URL.Query().Get("cursor")

	req := &list.Request{
		Limit:         limit,
		Order:         order,
		EncodedCursor: cursor,
	}

	res, err := h.list.Execute(ctx, req)
	if err != nil {
		expected := h.handleErr(ctx, w, err)
		if !expected {
			handleUnexpectedErr(ctx, w, err, "failed to list companies")
		}
		return
	}

	resp := mappers.CompanyListRespToDto(res)
	helpers.RespondJSON(ctx, w, http.StatusOK, resp)
}

// Create godoc
// @Summary Create new company
// @Description Creates new company, returns id
// @Tags company
// @Accept json
// @Produce json
// @Param company_request body dto.CompanyCreateRequest true "Request to create a company"
// @Success 201 {object} dto.CompanyCreatedResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /companies [post]
func (h *CompanyHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	iden := middleware.IdentityFromContext(ctx)
	dtoReq := middleware.BodyFromContext[dto.CompanyCreateRequest](ctx)

	req := mappers.CompanyCreateReqToUC(dtoReq)

	resp, err := h.create.Execute(ctx, req, iden)
	if err != nil {
		expected := h.handleErr(ctx, w, err)
		if !expected {
			handleUnexpectedErr(ctx, w, err, "failed to create company", "name", req.Name)
		}
		return
	}

	dtoResp := mappers.CompanyCreateRespToDto(resp)
	helpers.RespondJSON(ctx, w, http.StatusCreated, dtoResp)
}

// Update godoc
// @Summary Update company
// @Description Partially updates company fields, if field not provided or null - it won't be changed
// @Tags company
// @Accept json
// @Produce json
// @Param id path string true "Company ID"
// @Param company_request body dto.CompanyUpdateRequest true "Request to update company"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /companies/{id} [patch]
func (h *CompanyHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	iden := middleware.IdentityFromContext(ctx)
	id := middleware.UUIDFromContext(ctx, "id")
	dtoReq := middleware.BodyFromContext[dto.CompanyUpdateRequest](ctx)

	req := mappers.CompanyUpdateReqToUC(id, dtoReq)

	err := h.update.Execute(ctx, req, iden)
	if err != nil {
		expected := h.handleErr(ctx, w, err)
		if !expected {
			handleUnexpectedErr(ctx, w, err, "failed to update company", "id", id)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Delete godoc
// @Summary Delete company
// @Description Deletes company by id
// @Tags company
// @Produce json
// @Param id path string true "Company ID"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /companies/{id} [delete]
func (h *CompanyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	iden := middleware.IdentityFromContext(ctx)
	id := middleware.UUIDFromContext(ctx, "id")

	err := h.delete.Execute(ctx, id, iden)
	if err != nil {
		expected := h.handleErr(ctx, w, err)
		if !expected {
			handleUnexpectedErr(ctx, w, err, "failed to delete company", "id", id)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CompanyHandler) handleErr(ctx context.Context, w http.ResponseWriter, err error) bool {
	switch {
	case errors.Is(err, company.ErrCompanyNotFound):
		helpers.RespondError(ctx, w, http.StatusNotFound, err)
		return true

	case errors.Is(err, company.ErrCompanyAlreadyExists),
		errors.Is(err, member.ErrCompanyMemberAlreadyExists):
		helpers.RespondError(ctx, w, http.StatusConflict, err)
		return true

	case errors.Is(err, common.ErrInvalidCursor),
		errors.Is(err, common.ErrCursorOrderMismatch),
		errors.Is(err, company.ErrCompanyInvalidDescriptionLen),
		errors.Is(err, company.ErrCompanyInvalidNameLen),
		errors.Is(err, member.ErrInvalidUserID),
		errors.Is(err, member.ErrInvalidCompanyMemberRole):
		helpers.RespondError(ctx, w, http.StatusBadRequest, err)
		return true

	case errors.Is(err, identity.ErrInsufficientRole),
		errors.Is(err, identity.ErrHrRoleRequired),
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

func handleUnexpectedErr(ctx context.Context, w http.ResponseWriter, err error, ctxMsg string, logArgs ...any) {
	logger := utilslog.FromContext(ctx)
	logArgs = append(logArgs, "err", err)
	logger.ErrorContext(ctx, ctxMsg, logArgs...)
	helpers.RespondErrorMsg(ctx, w, http.StatusInternalServerError, "unexpected error")
}

func (h *CompanyHandler) parseOrderQuery(r *http.Request) list.Order {
	str := r.URL.Query().Get("order")
	ord := list.Order(strings.Trim(str, " "))

	switch ord {
	case list.OrderNameAsc,
		list.OrderCreatedAtDesc,
		list.OrderVacanciesDesc:
		return ord
	default:
		return list.OrderVacanciesDesc
	}
}
