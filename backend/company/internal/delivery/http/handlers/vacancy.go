package handlers

import (
	"errors"
	"net/http"

	"github.com/M0s1ck/g-store/src/pkg/http/middleware"
	"github.com/M0s1ck/g-store/src/pkg/http/responds"

	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/dto"
	_ "github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/helpers"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/mapper"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/create"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/get_by_id"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/update"
)

type VacancyHandler struct {
	getByID *get_vacancy.Usecase
	create  *create_vacancy.Usecase
	update  *update_vacancy.Usecase
}

func NewVacancyHandler(
	getByID *get_vacancy.Usecase,
	create *create_vacancy.Usecase,
	update *update_vacancy.Usecase,
) *VacancyHandler {

	return &VacancyHandler{
		getByID: getByID,
		create:  create,
		update:  update,
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

	if dtoReq.Title == "" {
		responds.RespondError(w, http.StatusBadRequest, errors.New("non empty title is required"))
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

func (h *VacancyHandler) handleErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain_errors.ErrVacancyNotFound):
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
		errors.Is(err, domain_errors.ErrInvalidDescriptionLength):
		responds.RespondError(w, http.StatusBadRequest, err)

	default:
		responds.RespondError(w, http.StatusInternalServerError, err)
	}
}
