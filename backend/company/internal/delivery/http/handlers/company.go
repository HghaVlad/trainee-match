package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/M0s1ck/g-store/src/pkg/http/middleware"
	"github.com/M0s1ck/g-store/src/pkg/http/responds"

	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/helpers"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/mapper"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/middleware"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/create"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/delete"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/get"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/list"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/update"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list_by_company"
)

type CompanyHandler struct {
	getByID *get_company.GetByIDUsecase
	create  *create_company.Usecase
	list    *list_companies.Usecase
	listVac *list_vac_by_comp.Usecase
	update  *update_company.Usecase
	delete  *delete_company.Usecase
}

func NewCompanyHandler(
	getByID *get_company.GetByIDUsecase,
	create *create_company.Usecase,
	list *list_companies.Usecase,
	update *update_company.Usecase,
	delete *delete_company.Usecase,
) *CompanyHandler {

	return &CompanyHandler{
		getByID: getByID,
		create:  create,
		list:    list,
		update:  update,
		delete:  delete,
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

	company, err := h.getByID.Execute(ctx, id)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	resp := mapper.GetCompRespToDto(company)
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

	req := &list_companies.Request{
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

	identity := my_middleware.IdentityFromContext(ctx)

	dtoReq, err := middleware.BodyFromContext[dto.CompanyCreateRequest](ctx)
	if err != nil {
		responds.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	req := mapper.CompanyCreateReqToUC(dtoReq)

	resp, err := h.create.Execute(ctx, req, identity)
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

	identity := my_middleware.IdentityFromContext(ctx)

	id, err := middleware.UUIDFromContext(ctx)
	if err != nil {
		responds.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	dtoReq, err := middleware.BodyFromContext[dto.CompanyUpdateRequest](ctx)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	req := mapper.CompanyUpdateReqToUC(id, dtoReq)

	err = h.update.Execute(ctx, req, identity)
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

	identity := my_middleware.IdentityFromContext(ctx)

	id, err := middleware.UUIDFromContext(ctx)
	if err != nil {
		responds.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.delete.Execute(ctx, id, identity)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CompanyHandler) handleErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain_errors.ErrCompanyNotFound):
		responds.RespondError(w, http.StatusNotFound, err)

	case errors.Is(err, domain_errors.ErrCompanyAlreadyExists),
		errors.Is(err, domain_errors.ErrCompanyMemberAlreadyExists):
		responds.RespondError(w, http.StatusConflict, err)

	case errors.Is(err, domain_errors.ErrInvalidCursor),
		errors.Is(err, domain_errors.ErrCursorOrderMismatch),
		errors.Is(err, domain_errors.ErrCompanyInvalidDescriptionLen),
		errors.Is(err, domain_errors.ErrCompanyInvalidNameLen),
		errors.Is(err, domain_errors.ErrInvalidUserID),
		errors.Is(err, domain_errors.ErrInvalidCompanyMemberRole):
		responds.RespondError(w, http.StatusBadRequest, err)

	case errors.Is(err, domain_errors.ErrInsufficientRole),
		errors.Is(err, domain_errors.ErrHrRoleRequired),
		errors.Is(err, domain_errors.ErrCompanyMemberRequired),
		errors.Is(err, domain_errors.ErrInsufficientRoleInCompany):
		responds.RespondError(w, http.StatusForbidden, err)

	default:
		responds.RespondError(w, http.StatusInternalServerError, err)
	}
}

func (h *CompanyHandler) parseOrderQuery(r *http.Request) list_companies.Order {
	str := r.URL.Query().Get("order")
	ord := list_companies.Order(strings.Trim(str, " "))

	switch ord {
	case list_companies.OrderNameAsc,
		list_companies.OrderCreatedAtDesc,
		list_companies.OrderVacanciesDesc:
		return ord
	default:
		return list_companies.OrderVacanciesDesc
	}
}
