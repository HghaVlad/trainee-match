package handlers

import (
	"errors"
	"net/http"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/M0s1ck/g-store/src/pkg/http/middleware"
	"github.com/M0s1ck/g-store/src/pkg/http/responds"

	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/mapper"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/get_company"
)

type CompanyHandler struct {
	getByID *get_company.GetByIDUsecase
}

func NewProfileHandler(
	getByID *get_company.GetByIDUsecase,

) *CompanyHandler {

	return &CompanyHandler{
		getByID: getByID,
	}
}

// GetById godoc
// @Summary Get profile by id
// @Description Returns company profile by UUID
// @Tags profile
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

	id, err := middleware.UUIDFromContext(ctx)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	// TODO: add timeout mb later

	company, err := h.getByID.Execute(ctx, id)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	resp := mapper.GetRespToDto(company)
	responds.RespondJSON(w, http.StatusOK, resp)
}

func (h *CompanyHandler) handleErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain_errors.ErrCompanyNotFound):
		responds.RespondError(w, http.StatusNotFound, err)

	default:
		responds.RespondError(w, http.StatusInternalServerError, err)
	}
}
