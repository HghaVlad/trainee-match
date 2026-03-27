package handlers

import (
	"errors"
	"net/http"
	"strings"

	gmiddleware "github.com/M0s1ck/g-store/src/pkg/http/middleware"
	"github.com/M0s1ck/g-store/src/pkg/http/responds"

	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/helpers"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/mapper"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/middleware"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/create"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/delete"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/get"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/list"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/update"
)

type CompanyHandler struct {
	getByID *get.Usecase
	create  *create.Usecase
	list    *list.Usecase
	update  *update.Usecase
	delete  *delete.Usecase
}

func NewCompanyHandler(
	get *get.Usecase,
	create *create.Usecase,
	list *list.Usecase,
	upd *update.Usecase,
	del *delete.Usecase,
) *CompanyHandler {

	return &CompanyHandler{
		getByID: get,
		create:  create,
		list:    list,
		update:  upd,
		delete:  del,
	}
}

// GetById godoc
// @Summary Get profile by id
// @Description Returns company profile by UUID
// @Tags company
// @Accept json
// @Produce json
// @Param id path string true "Company ID (UUID)"
// @Success 200 {object} dto.CompanyResponse
// @Failure 400 {object} responds.ErrorResponse
// @Failure 404 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /companies/{id} [get]
func (h *CompanyHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, ok := helpers.ParseUuidFromPathOr400(r, w, "id")
	if !ok {
		return
	}

	comp, err := h.getByID.Execute(ctx, id)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	resp := mapper.GetCompRespToDto(comp)
	responds.RespondJSON(w, http.StatusOK, resp)
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
// @Failure 400 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /companies [get]
func (h *CompanyHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := helpers.ParseLimit(r, "limit", 20)
	order := h.parseOrderQuery(r)
	cursor := r.URL.Query().Get("cursor")

	req := &list.Request{
		Limit:  limit,
		Order:  order,
		Cursor: cursor,
	}

	res, err := h.list.Execute(ctx, req)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	resp := mapper.CompanyListRespToDto(res)
	responds.RespondJSON(w, http.StatusOK, resp)
}

// Create godoc
// @Summary Create new company
// @Description Creates new company, returns id
// @Tags company
// @Accept json
// @Produce json
// @Param company_request body dto.CompanyCreateRequest true "Request to create a company"
// @Success 201 {object} dto.CompanyCreatedResponse
// @Failure 400 {object} responds.ErrorResponse
// @Failure 401 {object} responds.ErrorResponse
// @Failure 403 {object} responds.ErrorResponse
// @Failure 409 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /companies [post]
func (h *CompanyHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	iden := middleware.IdentityFromContext(ctx)

	dtoReq, err := gmiddleware.BodyFromContext[dto.CompanyCreateRequest](ctx)
	if err != nil {
		responds.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	req := mapper.CompanyCreateReqToUC(dtoReq)

	resp, err := h.create.Execute(ctx, req, iden)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	dtoResp := mapper.CompanyCreateRespToDto(resp)
	responds.RespondJSON(w, http.StatusCreated, dtoResp)
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
// @Failure 400 {object} responds.ErrorResponse
// @Failure 401 {object} responds.ErrorResponse
// @Failure 403 {object} responds.ErrorResponse
// @Failure 404 {object} responds.ErrorResponse
// @Failure 409 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /companies/{id} [patch]
func (h *CompanyHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	iden := middleware.IdentityFromContext(ctx)

	id, err := gmiddleware.UUIDFromContext(ctx)
	if err != nil {
		responds.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	dtoReq, err := gmiddleware.BodyFromContext[dto.CompanyUpdateRequest](ctx)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	req := mapper.CompanyUpdateReqToUC(id, dtoReq)

	err = h.update.Execute(ctx, req, iden)
	if err != nil {
		h.handleErr(w, err)
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
// @Failure 400 {object} responds.ErrorResponse
// @Failure 401 {object} responds.ErrorResponse
// @Failure 403 {object} responds.ErrorResponse
// @Failure 404 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /companies/{id} [delete]
func (h *CompanyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	iden := middleware.IdentityFromContext(ctx)

	id, err := gmiddleware.UUIDFromContext(ctx)
	if err != nil {
		responds.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.delete.Execute(ctx, id, iden)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CompanyHandler) handleErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, company.ErrCompanyNotFound):
		responds.RespondError(w, http.StatusNotFound, err)

	case errors.Is(err, company.ErrCompanyAlreadyExists),
		errors.Is(err, member.ErrCompanyMemberAlreadyExists):
		responds.RespondError(w, http.StatusConflict, err)

	case errors.Is(err, common.ErrInvalidCursor),
		errors.Is(err, common.ErrCursorOrderMismatch),
		errors.Is(err, company.ErrCompanyInvalidDescriptionLen),
		errors.Is(err, company.ErrCompanyInvalidNameLen),
		errors.Is(err, member.ErrInvalidUserID),
		errors.Is(err, member.ErrInvalidCompanyMemberRole):
		responds.RespondError(w, http.StatusBadRequest, err)

	case errors.Is(err, identity.ErrInsufficientRole),
		errors.Is(err, identity.ErrHrRoleRequired),
		errors.Is(err, member.ErrCompanyMemberRequired),
		errors.Is(err, member.ErrInsufficientRoleInCompany):
		responds.RespondError(w, http.StatusForbidden, err)

	default:
		responds.RespondError(w, http.StatusInternalServerError, err)
	}
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
