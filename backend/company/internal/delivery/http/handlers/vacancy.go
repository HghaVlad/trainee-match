package handlers

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/archive"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/publish"
	"github.com/M0s1ck/g-store/src/pkg/http/middleware"
	"github.com/M0s1ck/g-store/src/pkg/http/responds"
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/helpers"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/mapper"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/middleware"
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
	publish    *publish_vacancy.Usecase
	archive    *archive_vacancy.Usecase
	delete     *delete_vacancy.Usecase
}

func NewVacancyHandler(
	getByID *get_vacancy.Usecase,
	list *list_vacancy.Usecase,
	listByComp *list_vac_by_comp.Usecase,
	create *create_vacancy.Usecase,
	update *update_vacancy.Usecase,
	publish *publish_vacancy.Usecase,
	archive *archive_vacancy.Usecase,
	delete *delete_vacancy.Usecase,
) *VacancyHandler {

	return &VacancyHandler{
		getByID:    getByID,
		list:       list,
		listByComp: listByComp,
		create:     create,
		update:     update,
		publish:    publish,
		archive:    archive,
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
// @Success 200 {object} dto.VacancyFullResponse
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
// @Failure 401 {object} responds.ErrorResponse
// @Failure 403 {object} responds.ErrorResponse
// @Failure 404 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /companies/{company-id}/vacancies [post]
func (h *VacancyHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	identity := my_middleware.IdentityFromContext(ctx)

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

	resp, err := h.create.Execute(ctx, req, identity)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	dtoResp := mapper.VacancyCreateRespToDto(resp)
	responds.RespondJSON(w, http.StatusCreated, dtoResp)
}

// List godoc
// @Summary List vacancy summaries
// @Description Uses cursor pagination, returns next cursor if there's more. Supports filters, orders.
// @Tags vacancy
// @Accept json
// @Produce json
// @Param order query string false "Order attribute, supports published_at_desc, salary_desc, salary_asc" default(published_at_desc)
// @Param cursor query string false "Cursor"
// @Param limit query int false "Items per page" default(20)
// / Filters:
// Salary range
// @Param salary_min query int false "Minimum salary"
// @Param salary_max query int false "Maximum salary"
// Hours per week range
// @Param hours_min query int false "Minimum hours per week"
// @Param hours_max query int false "Maximum hours per week"
// Duration range (months)
// @Param duration_min query int false "Minimum duration in days"
// @Param duration_max query int false "Maximum duration in days"
// Boolean filters
// @Param is_paid query bool false "Paid vacancy filter"
// @Param internship_to_offer query bool false "Internship with possible job offer"
// @Param flexible_schedule query bool false "Flexible schedule filter"
// Slice filters (multiple values allowed)
// @Param work_format query []string false "Work format filter (repeat param)" collectionFormat(multi)
// @Param city query []string false "City filter (repeat param)" collectionFormat(multi)
// @Param company_id query []string false "Company filter (repeat param)" collectionFormat(multi)
// @Success 200 {object} dto.VacancyListResponse
// @Failure 400 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /vacancies [get]
func (h *VacancyHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := h.listVacRequestFromQuery(r)
	if err != nil {
		responds.RespondError(w, http.StatusBadRequest, err)
		return
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
// @Summary Lists company's vacancy summaries. Outdated, needs update if needed. Rn u can use list with company_id param
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
// @Failure 401 {object} responds.ErrorResponse
// @Failure 403 {object} responds.ErrorResponse
// @Failure 404 {object} responds.ErrorResponse
// @Failure 409 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /companies/{company-id}/vacancies/{vacancy-id} [patch]
func (h *VacancyHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	identity := my_middleware.IdentityFromContext(ctx)

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

	err = h.update.Execute(ctx, req, identity)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Publish
// @Summary Publish vacancy
// @Description Publish vacancy for candidates
// @Tags vacancy
// @Accept json
// @Produce json
// @Param company-id path string true "Company ID"
// @Param vacancy-id path string true "Vacancy ID"
// @Success 204 "Vacancy published successfully"
// @Failure 400 {object} responds.ErrorResponse
// @Failure 401 {object} responds.ErrorResponse
// @Failure 403 {object} responds.ErrorResponse
// @Failure 404 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /companies/{company-id}/vacancies/{vacancy-id}/publish [post]
func (h *VacancyHandler) Publish(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	identity := my_middleware.IdentityFromContext(ctx)

	companyID, ok := helpers.ParseUuidFromPathOr400(r, w, "company-id")
	if !ok {
		return
	}

	vacancyID, ok := helpers.ParseUuidFromPathOr400(r, w, "vacancy-id")
	if !ok {
		return
	}

	err := h.publish.Execute(ctx, companyID, vacancyID, identity)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Archive
// @Summary Archive vacancy (deactivation for candidates)
// @Description Archive vacancy (deactivation for candidates)
// @Tags vacancy
// @Accept json
// @Produce json
// @Param company-id path string true "Company ID"
// @Param vacancy-id path string true "Vacancy ID"
// @Success 204 "Vacancy archived successfully"
// @Failure 400 {object} responds.ErrorResponse
// @Failure 401 {object} responds.ErrorResponse
// @Failure 403 {object} responds.ErrorResponse
// @Failure 404 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /companies/{company-id}/vacancies/{vacancy-id}/archive [post]
func (h *VacancyHandler) Archive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	identity := my_middleware.IdentityFromContext(ctx)

	companyID, ok := helpers.ParseUuidFromPathOr400(r, w, "company-id")
	if !ok {
		return
	}

	vacancyID, ok := helpers.ParseUuidFromPathOr400(r, w, "vacancy-id")
	if !ok {
		return
	}

	err := h.archive.Execute(ctx, companyID, vacancyID, identity)
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
// @Failure 401 {object} responds.ErrorResponse
// @Failure 403 {object} responds.ErrorResponse
// @Failure 404 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /companies/{company-id}/vacancies/{vacancy-id} [delete]
func (h *VacancyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	identity := my_middleware.IdentityFromContext(ctx)

	companyID, ok := helpers.ParseUuidFromPathOr400(r, w, "company-id")
	if !ok {
		return
	}

	vacancyID, ok := helpers.ParseUuidFromPathOr400(r, w, "vacancy-id")
	if !ok {
		return
	}

	err := h.delete.Execute(ctx, vacancyID, companyID, identity)
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
		errors.Is(err, domain_errors.ErrInvalidCursor),
		errors.Is(err, domain_errors.ErrCursorOrderMismatch),
		errors.Is(err, domain_errors.ErrUnsupportedListOrder),
		errors.Is(err, domain_errors.ErrEmptyCityFilter),
		errors.Is(err, domain_errors.ErrEmptyCompaniesFilter),
		errors.Is(err, domain_errors.ErrInvalidSalaryOrderForUnpaid),
		errors.Is(err, domain_errors.ErrLimitTooLarge):
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

func (h *VacancyHandler) listVacRequestFromQuery(r *http.Request) (*list_vacancy.Request, error) {
	q := r.URL.Query()

	limit := helpers.ParseLimit(r, "limit", 20)
	order, err := h.parseListOrderQuery(r)
	if err != nil {
		return nil, err
	}
	cursor := q.Get("cursor")

	req := &list_vacancy.Request{
		Limit:         limit,
		Order:         order,
		EncodedCursor: cursor,
		Requirements:  new(list_vacancy.Requirements),
	}

	req.Requirements.Salary = parseRangeInt(q, "salary_min", "salary_max")
	req.Requirements.HoursPerWeek = parseRangeInt(q, "hours_min", "hours_max")
	req.Requirements.Duration = parseRangeInt(q, "duration_min", "duration_max")

	if isPaidStr := q.Get("is_paid"); isPaidStr != "" {
		if isPaid, err := strconv.ParseBool(isPaidStr); err == nil {
			req.Requirements.IsPaid = &isPaid
		}
	}

	if internshipStr := q.Get("internship_to_offer"); internshipStr != "" {
		if val, err := strconv.ParseBool(internshipStr); err == nil {
			req.Requirements.InternshipToOffer = &val
		}
	}

	if flexStr := q.Get("flexible_schedule"); flexStr != "" {
		if val, err := strconv.ParseBool(flexStr); err == nil {
			req.Requirements.FlexibleSchedule = &val
		}
	}

	if workFormats, ok := q["work_format"]; ok && len(workFormats) > 0 {
		var wfs []value_types.WorkFormat
		for _, str := range workFormats {
			wf := value_types.WorkFormat(str)
			if wf.IsValid() {
				wfs = append(wfs, wf)
			}
		}
		if len(wfs) > 0 {
			req.Requirements.WorkFormat = &wfs
		}
	}

	if companies, ok := q["company_id"]; ok && len(companies) > 0 {
		ids := make([]uuid.UUID, 0, len(companies))
		for _, str := range companies {
			id, err := uuid.Parse(str)
			if err == nil {
				ids = append(ids, id)
			}
		}
		req.Requirements.Companies = &ids
	}

	if cities, ok := q["city"]; ok && len(cities) > 0 {
		req.Requirements.City = &cities
	}

	return req, nil
}

func parseRangeInt(q url.Values, minKey, maxKey string) *list_vacancy.RangeInt {
	var r list_vacancy.RangeInt
	var hasValue bool

	if minStr := q.Get(minKey); minStr != "" {
		if mn, err := strconv.Atoi(minStr); err == nil {
			r.Min = &mn
			hasValue = true
		}
	}

	if maxStr := q.Get(maxKey); maxStr != "" {
		if mx, err := strconv.Atoi(maxStr); err == nil {
			r.Max = &mx
			hasValue = true
		}
	}

	if !hasValue {
		return nil
	}

	return &r
}

func (h *VacancyHandler) parseListOrderQuery(r *http.Request) (list_vacancy.Order, error) {
	str := r.URL.Query().Get("order")
	if str == "" {
		return list_vacancy.OrderPublishedAtDesc, nil
	}

	ord := list_vacancy.Order(strings.Trim(str, " "))

	switch ord {
	case list_vacancy.OrderPublishedAtDesc,
		list_vacancy.OrderSalaryDesc,
		list_vacancy.OrderSalaryAsc:
		return ord, nil
	default:
		return "", domain_errors.ErrUnsupportedListOrder
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
