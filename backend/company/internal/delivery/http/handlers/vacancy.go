package handlers

import (
	"errors"
	"net/http"

	"github.com/M0s1ck/g-store/src/pkg/http/responds"

	_ "github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/helpers"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/mapper"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/get_by_id"
)

type VacancyHandler struct {
	getByID *get_vacancy.Usecase
}

func NewVacancyHandler(
	getByID *get_vacancy.Usecase,
) *VacancyHandler {

	return &VacancyHandler{
		getByID: getByID,
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

func (h *VacancyHandler) handleErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain_errors.ErrVacancyNotFound):
		responds.RespondError(w, http.StatusNotFound, err)

	default:
		responds.RespondError(w, http.StatusInternalServerError, err)
	}
}
