package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/dto"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/helpers"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/mappers"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/middleware"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/archive"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/create"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/get"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/getpublished"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/listbycomp"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/publish"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/remove"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/update"
)

type VacancyHandler struct {
	getByID          *get.Usecase
	getPublishedByID *getpublished.Usecase
	list             *list.Usecase
	listByComp       *listbycomp.Usecase
	create           *create.Usecase
	update           *update.Usecase
	publish          *publish.Usecase
	archive          *archive.Usecase
	del              *remove.Usecase
}

func NewVacancyHandler(
	getByID *get.Usecase,
	getPublishedByID *getpublished.Usecase,
	list *list.Usecase,
	listByComp *listbycomp.Usecase,
	create *create.Usecase,
	update *update.Usecase,
	publish *publish.Usecase,
	archive *archive.Usecase,
	del *remove.Usecase,
) *VacancyHandler {
	return &VacancyHandler{
		getByID:          getByID,
		getPublishedByID: getPublishedByID,
		list:             list,
		listByComp:       listByComp,
		create:           create,
		update:           update,
		publish:          publish,
		archive:          archive,
		del:              del,
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
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /companies/{company-id}/vacancies/{vacancy-id} [get]
func (h *VacancyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	iden := middleware.IdentityFromContext(ctx)
	companyID := middleware.UUIDFromContext(ctx, "company-id")
	vacancyID := middleware.UUIDFromContext(ctx, "vacancy-id")

	vac, err := h.getByID.Execute(ctx, vacancyID, companyID, iden)
	if err != nil {
		expected := h.handleErr(ctx, w, err)
		if !expected {
			handleUnexpectedErr(ctx, w, err, "failed to get vacancy",
				"id", vacancyID)
		}
		return
	}

	resp := mappers.VacancyToDtoResponse(vac)
	helpers.RespondJSON(ctx, w, http.StatusOK, resp)
}

// GetPublishedByID godoc
// @Summary Get published vacancy by id
// @Description Returns public vacancy view for candidates. Only published vacancies are visible.
// @Tags vacancy
// @Accept json
// @Produce json
// @Param vacancy-id path string true "Vacancy ID (UUID)"
// @Success 200 {object} dto.VacancyPublicResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /vacancies/{vacancy-id} [get]
func (h *VacancyHandler) GetPublishedByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vacancyID := middleware.UUIDFromContext(ctx, "id")

	vac, err := h.getPublishedByID.Execute(ctx, vacancyID)
	if err != nil {
		expected := h.handleErr(ctx, w, err)
		if !expected {
			handleUnexpectedErr(ctx, w, err, "failed to get vacancy",
				"id", vacancyID)
		}
		return
	}

	resp := mappers.VacancyPublicToDtoResponse(vac)
	helpers.RespondJSON(ctx, w, http.StatusOK, resp)
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
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /companies/{company-id}/vacancies [post]
func (h *VacancyHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	iden := middleware.IdentityFromContext(ctx)
	companyID := middleware.UUIDFromContext(ctx, "company-id")
	dtoReq := middleware.BodyFromContext[dto.VacancyCreateRequest](ctx)

	req := mappers.VacancyCreateReqToUC(dtoReq, companyID)

	resp, err := h.create.Execute(ctx, req, iden)
	if err != nil {
		expected := h.handleErr(ctx, w, err)
		if !expected {
			handleUnexpectedErr(ctx, w, err, "failed to create vacancy",
				"company_id", companyID)
		}
		return
	}

	dtoResp := mappers.VacancyCreateRespToDto(resp)
	helpers.RespondJSON(ctx, w, http.StatusCreated, dtoResp)
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
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /vacancies [get]
func (h *VacancyHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := helpers.ListVacRequestFromQuery(r)
	if err != nil {
		helpers.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	res, err := h.list.Execute(ctx, req)
	if err != nil {
		expected := h.handleErr(ctx, w, err)
		if !expected {
			handleUnexpectedErr(ctx, w, err, "failed to list vacancy")
		}
		return
	}

	resp := mappers.VacancyListRespToDto(res)
	helpers.RespondJSON(ctx, w, http.StatusOK, resp)
}

// ListByCompany godoc
// @Summary Lists company's vacancy summaries
// @Description Uses cursor pagination, returns next cursor if there's more. Supports filters and status.
// @Tags vacancy
// @Accept json
// @Produce json
// @Param company-id path string true "Company ID (UUID)"
// @Param order query string false "Order attribute" default(created_at_desc)
// @Param cursor query string false "Cursor"
// @Param limit query int false "Items per page" default(20)
// @Param status query string false "Vacancy status filter"
// @Param salary_min query int false "Minimum salary"
// @Param salary_max query int false "Maximum salary"
// @Param hours_min query int false "Minimum hours per week"
// @Param hours_max query int false "Maximum hours per week"
// @Param duration_min query int false "Minimum duration in days"
// @Param duration_max query int false "Maximum duration in days"
// @Param is_paid query bool false "Paid vacancy filter"
// @Param internship_to_offer query bool false "Internship with possible job offer"
// @Param flexible_schedule query bool false "Flexible schedule filter"
// @Param work_format query []string false "Work format filter (repeat param)" collectionFormat(multi)
// @Param city query []string false "City filter (repeat param)" collectionFormat(multi)
// @Success 200 {object} dto.VacancyByCompListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /companies/{company-id}/vacancies [get]
func (h *VacancyHandler) ListByCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	iden := middleware.IdentityFromContext(ctx)
	compID := middleware.UUIDFromContext(ctx, "company-id")

	req, err := helpers.ListVacByCompRequestFromQuery(r, compID)
	if err != nil {
		helpers.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	res, err := h.listByComp.Execute(ctx, req, iden)
	if err != nil {
		expected := h.handleErr(ctx, w, err)
		if !expected {
			handleUnexpectedErr(ctx, w, err, "failed to list vacancy by company")
		}
		return
	}

	resp := mappers.ListVacByCompRespToDto(res)
	helpers.RespondJSON(ctx, w, http.StatusOK, resp)
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
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /companies/{company-id}/vacancies/{vacancy-id} [patch]
func (h *VacancyHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	iden := middleware.IdentityFromContext(ctx)
	companyID := middleware.UUIDFromContext(ctx, "company-id")
	vacancyID := middleware.UUIDFromContext(ctx, "vacancy-id")
	dtoReq := middleware.BodyFromContext[dto.VacancyUpdateRequest](ctx)

	req := mappers.VacancyUpdateReqToUC(dtoReq, companyID, vacancyID)

	err := h.update.Execute(ctx, req, iden)
	if err != nil {
		expected := h.handleErr(ctx, w, err)
		if !expected {
			handleUnexpectedErr(ctx, w, err, "failed to update vacancy",
				"id", vacancyID)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Publish godoc
// @Summary Publish vacancy
// @Description Publish vacancy for candidates
// @Tags vacancy
// @Accept json
// @Produce json
// @Param company-id path string true "Company ID"
// @Param vacancy-id path string true "Vacancy ID"
// @Success 204 "Vacancy published successfully"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /companies/{company-id}/vacancies/{vacancy-id}/publish [post]
func (h *VacancyHandler) Publish(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	iden := middleware.IdentityFromContext(ctx)
	companyID := middleware.UUIDFromContext(ctx, "company-id")
	vacancyID := middleware.UUIDFromContext(ctx, "vacancy-id")

	err := h.publish.Execute(ctx, companyID, vacancyID, iden)
	if err != nil {
		expected := h.handleErr(ctx, w, err)
		if !expected {
			handleUnexpectedErr(ctx, w, err, "failed to publish vacancy",
				"id", vacancyID)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Archive godoc
// @Summary Archive vacancy (deactivation for candidates)
// @Description Archive vacancy (deactivation for candidates)
// @Tags vacancy
// @Accept json
// @Produce json
// @Param company-id path string true "Company ID"
// @Param vacancy-id path string true "Vacancy ID"
// @Success 204 "Vacancy archived successfully"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /companies/{company-id}/vacancies/{vacancy-id}/archive [post]
func (h *VacancyHandler) Archive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	iden := middleware.IdentityFromContext(ctx)
	companyID := middleware.UUIDFromContext(ctx, "company-id")
	vacancyID := middleware.UUIDFromContext(ctx, "vacancy-id")

	err := h.archive.Execute(ctx, companyID, vacancyID, iden)
	if err != nil {
		expected := h.handleErr(ctx, w, err)
		if !expected {
			handleUnexpectedErr(ctx, w, err, "failed to archive vacancy",
				"id", vacancyID)
		}
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
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /companies/{company-id}/vacancies/{vacancy-id} [delete]
func (h *VacancyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	iden := middleware.IdentityFromContext(ctx)
	companyID := middleware.UUIDFromContext(ctx, "company-id")
	vacancyID := middleware.UUIDFromContext(ctx, "vacancy-id")

	err := h.del.Execute(ctx, vacancyID, companyID, iden)
	if err != nil {
		expected := h.handleErr(ctx, w, err)
		if !expected {
			handleUnexpectedErr(ctx, w, err, "failed to delete vacancy",
				"id", vacancyID)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *VacancyHandler) handleErr(ctx context.Context, w http.ResponseWriter, err error) bool {
	switch {
	case errors.Is(err, vacancy.ErrVacancyNotFound),
		errors.Is(err, company.ErrCompanyNotFound):
		helpers.RespondError(ctx, w, http.StatusNotFound, err)
		return true

	case errors.Is(err, vacancy.ErrInvalidWorkFormat),
		errors.Is(err, vacancy.ErrInvalidEmploymentType),
		errors.Is(err, vacancy.ErrInvalidDurationRange),
		errors.Is(err, vacancy.ErrInvalidHoursRange),
		errors.Is(err, vacancy.ErrInvalidSalaryRange),
		errors.Is(err, vacancy.ErrSalaryProvidedForUnpaid),
		errors.Is(err, vacancy.ErrSalaryMissingForPaid),
		errors.Is(err, vacancy.ErrNegativeSalary),
		errors.Is(err, vacancy.ErrSalaryTooLarge),
		errors.Is(err, vacancy.ErrInvalidTitleLength),
		errors.Is(err, vacancy.ErrInvalidDescriptionLength),
		errors.Is(err, vacancy.ErrInvalidStatus),
		errors.Is(err, vacancy.ErrEmptyCityFilter),
		errors.Is(err, vacancy.ErrEmptyCompaniesFilter),
		errors.Is(err, vacancy.ErrInvalidSalaryOrderForUnpaid),
		errors.Is(err, common.ErrInvalidCursor),
		errors.Is(err, common.ErrCursorOrderMismatch),
		errors.Is(err, common.ErrUnsupportedListOrder),
		errors.Is(err, common.ErrLimitTooLarge):
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
