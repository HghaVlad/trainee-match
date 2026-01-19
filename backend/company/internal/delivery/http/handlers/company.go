package handlers

import (
	"net/http"

	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/helpers"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/mapper"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/middleware"
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
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /companies/{id} [get]
func (h *CompanyHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := my_middleware.UUIDFromContext(ctx)

	// TODO: add timeout mb later

	company, err := h.getByID.Execute(ctx, id)
	if err != nil {
		helpers.HandleError(w, err)
		return
	}

	resp := mapper.GetRespToDto(company)
	helpers.RespondJSON(w, http.StatusOK, resp)
}
