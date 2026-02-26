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
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/create"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/delete"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/get_by_id"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list_by_company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/update"
)

type VacancyHandler struct {
	getByID    *get_vacancy.Usecase
	list       *list_vacancy.Usecase
	listByComp *list_vac_by_comp.Usecase
	create     *create_vacancy.Usecase
	update     *update_vacancy.Usecase
	delete     *delete_vacancy.Usecase
}

func NewVacancyHandler(
	getByID *get_vacancy.Usecase,
	list *list_vacancy.Usecase,
	listByComp *list_vac_by_comp.Usecase,
	create *create_vacancy.Usecase,
	update *update_vacancy.Usecase,
	delete *delete_vacancy.Usecase,
) *VacancyHandler {

	return &VacancyHandler{
		getByID:    getByID,
		list:       list,
		listByComp: listByComp,
		create:     create,
		update:     update,
		delete:     delete,
	}
}

// GetByID godoc
// @Summary Get vacancy by id
// @Description Returns vacancy by id, company id should be correct, otherwise it's 404
// @Tags vacancy
// @Accept json
// @Produce json
// @Param company-id path string true "Company ID (UUID)"
// @Param vacancy-id path string true "Vacancy ID (UUID)"
// @Success 200 {object} dto.VacancyResponse
// @Failure 400 {object} responds.ErrorResponse
// @Failure 404 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /companies/{company-id}/vacancies/{vacancy-id} [get]
func (h *VacancyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	companyID, ok := helpers.ParseUuidFromPathOr400(r, w, "company-id")
	if !ok {
		return
	}

	vacancyID, ok := helpers.ParseUuidFromPathOr400(r, w, "vacancy-id")
	if !ok {
		return
	}

	vacancy, err := h.getByID.Execute(ctx, vacancyID, companyID)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	resp := mapper.VacancyToDtoResponse(vacancy)
	responds.RespondJSON(w, http.StatusOK, resp)
}

// Create godoc
// @Summary Create new vacancy
// @Description Creates new vacancy, returns id
// @Tags vacancy
// @Accept json
// @Produce json
// @Param company-id path string true "Company ID (UUID)"
// @Param vacancy_request body dto.VacancyCreateRequest true "Request to create vacancy"
// @Success 201 {object} dto.VacancyCreatedResponse
// @Failure 400 {object} responds.ErrorResponse
// @Failure 404 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /companies/{company-id}/vacancies [post]
func (h *VacancyHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoReq, err := middleware.BodyFromContext[dto.VacancyCreateRequest](ctx)
	if err != nil {
		responds.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	companyID, ok := helpers.ParseUuidFromPathOr400(r, w, "company-id")
	if !ok {
		return
	}

	req := mapper.VacancyCreateReqToUC(dtoReq, companyID)

	resp, err := h.create.Execute(ctx, req)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	dtoResp := mapper.VacancyCreateRespToDto(resp)
	responds.RespondJSON(w, http.StatusCreated, dtoResp)
}

// List godoc
// @Summary List vacancy summaries
// @Description Uses cursor pagination, returns next cursor if there's more. Supports order by published_at_desc
// @Tags vacancy
// @Accept json
// @Produce json
// @Param order query string false "Order attribute" default(published_at_desc)
// @Param cursor query string false "Cursor"
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.VacancyListResponse
// @Failure 400 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /vacancies [get]
func (h *VacancyHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit := helpers.ParseLimit(r, "limit", 20)
	order := h.parseListOrderQuery(r)
	cursor := r.URL.Query().Get("cursor")

	req := &list_vacancy.Request{
		Limit:  limit,
		Order:  order,
		Cursor: cursor,
	}

	res, err := h.list.Execute(ctx, req)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	resp := mapper.VacancyListRespToDto(res)
	responds.RespondJSON(w, http.StatusOK, resp)
}

// ListByCompany godoc
// @Summary Lists company's vacancy summaries
// @Description Uses cursor pagination, returns next cursor if there's more. Supports order by published_at_desc
// @Tags vacancy
// @Accept json
// @Produce json
// @Param company-id path string true "Company ID (UUID)"
// @Param order query string false "Order attribute" default(published_at_desc)
// @Param cursor query string false "Cursor"
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.VacancyByCompListResponse
// @Failure 400 {object} responds.ErrorResponse
// @Failure 404 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /companies/{company-id}/vacancies [get]
func (h *VacancyHandler) ListByCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	compID, ok := helpers.ParseUuidFromPathOr400(r, w, "company-id")
	if !ok {
		return
	}

	order := h.parseVacByCompListOrderQuery(r)
	cursor := r.URL.Query().Get("cursor")
	limit := helpers.ParseLimit(r, "limit", 20)

	req := &list_vac_by_comp.Request{
		CompID: compID,
		Limit:  limit,
		Order:  order,
		Cursor: cursor,
	}

	res, err := h.listByComp.Execute(ctx, req)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	resp := mapper.ListVacByCompRespToDto(res)
	responds.RespondJSON(w, http.StatusOK, resp)
}

// Update godoc
// @Summary Update vacancy
// @Description Partially updates vacancy fields. Nil fields are ignored (not updated).
// @Tags vacancy
// @Accept json
// @Produce json
// @Param company-id path string true "Company ID"
// @Param vacancy-id path string true "Vacancy ID"
// @Param vacancy_request body dto.VacancyUpdateRequest true "Vacancy update payload"
// @Success 204 "Vacancy updated successfully"
// @Failure 400 {object} responds.ErrorResponse
// @Failure 404 {object} responds.ErrorResponse
// @Failure 409 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /companies/{company-id}/vacancies/{vacancy-id} [patch]
func (h *VacancyHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	companyID, ok := helpers.ParseUuidFromPathOr400(r, w, "company-id")
	if !ok {
		return
	}

	vacancyID, ok := helpers.ParseUuidFromPathOr400(r, w, "vacancy-id")
	if !ok {
		return
	}

	dtoReq, err := middleware.BodyFromContext[dto.VacancyUpdateRequest](ctx)
	if err != nil {
		responds.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	req := mapper.VacancyUpdateReqToUC(dtoReq, companyID, vacancyID)

	err = h.update.Execute(ctx, req)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Delete godoc
// @Summary Delete vacancy
// @Description Deletes vacancy by id
// @Tags vacancy
// @Produce json
// @Param vacancy-id path string true "Vacancy ID"
// @Param company-id path string true "Company ID"
// @Success 204
// @Failure 400 {object} responds.ErrorResponse
// @Failure 404 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /companies/{company-id}/vacancies/{vacancy-id} [delete]
func (h *VacancyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	companyID, ok := helpers.ParseUuidFromPathOr400(r, w, "company-id")
	if !ok {
		return
	}

	vacancyID, ok := helpers.ParseUuidFromPathOr400(r, w, "vacancy-id")
	if !ok {
		return
	}

	err := h.delete.Execute(ctx, vacancyID, companyID)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *VacancyHandler) handleErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain_errors.ErrVacancyNotFound),
		errors.Is(err, domain_errors.ErrCompanyNotFound):
		responds.RespondError(w, http.StatusNotFound, err)

	case errors.Is(err, domain_errors.ErrInvalidWorkFormat),
		errors.Is(err, domain_errors.ErrInvalidEmploymentType),
		errors.Is(err, domain_errors.ErrInvalidDurationRange),
		errors.Is(err, domain_errors.ErrInvalidHoursRange),
		errors.Is(err, domain_errors.ErrInvalidSalaryRange),
		errors.Is(err, domain_errors.ErrSalaryProvidedForUnpaid),
		errors.Is(err, domain_errors.ErrSalaryMissingForPaid),
		errors.Is(err, domain_errors.ErrNegativeSalary),
		errors.Is(err, domain_errors.ErrSalaryTooLarge),
		errors.Is(err, domain_errors.ErrInvalidTitleLength),
		errors.Is(err, domain_errors.ErrInvalidDescriptionLength),
		errors.Is(err, domain_errors.ErrInvalidCursor):
		responds.RespondError(w, http.StatusBadRequest, err)

	default:
		responds.RespondError(w, http.StatusInternalServerError, err)
	}
}

func (h *VacancyHandler) parseListOrderQuery(r *http.Request) list_vacancy.Order {
	str := r.URL.Query().Get("order")
	ord := list_vacancy.Order(strings.Trim(str, " "))

	switch ord {
	case list_vacancy.OrderPublishedAtDesc:
		return ord
	default:
		return list_vacancy.OrderPublishedAtDesc
	}
}

func (h *VacancyHandler) parseVacByCompListOrderQuery(r *http.Request) list_vac_by_comp.Order {
	str := r.URL.Query().Get("order")
	ord := list_vac_by_comp.Order(strings.Trim(str, " "))

	switch ord {
	case list_vac_by_comp.OrderPublishedAtDesc:
		return ord
	default:
		return list_vac_by_comp.OrderPublishedAtDesc
	}
}
