package handlers

import (
	"errors"
	"net/http"

	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/mapper"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/create_company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/delete_company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/get_company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/update_company"
	"github.com/M0s1ck/g-store/src/pkg/http/middleware"
	"github.com/M0s1ck/g-store/src/pkg/http/responds"
)

type CompanyHandler struct {
	getByID *get_company.GetByIDUsecase
	create  *create_company.Usecase
	update  *update_company.Usecase
	delete  *delete_company.Usecase
}

func NewProfileHandler(
	getByID *get_company.GetByIDUsecase,
	create *create_company.Usecase,
	update *update_company.Usecase,
	delete *delete_company.Usecase,
) *CompanyHandler {

	return &CompanyHandler{
		getByID: getByID,
		create:  create,
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

	resp := mapper.GetCompRespToDto(company)
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
// @Failure 409 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /companies [post]
func (h *CompanyHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoReq, err := middleware.BodyFromContext[dto.CompanyCreateRequest](ctx)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	if dtoReq.Name == "" {
		responds.RespondError(w, http.StatusBadRequest, errors.New("non empty name is required"))
		return
	}

	// TODO: add timeout mb later
	// TODO: add jwt owner id prolly

	req := mapper.CompanyCreateReqToUC(dtoReq)

	resp, err := h.create.Execute(ctx, req)
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
// @Failure 404 {object} responds.ErrorResponse
// @Failure 409 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /companies/{id} [patch]
func (h *CompanyHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := middleware.UUIDFromContext(ctx)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	dtoReq, err := middleware.BodyFromContext[dto.CompanyUpdateRequest](ctx)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	if dtoReq.Name != nil && *dtoReq.Name == "" {
		responds.RespondError(w, http.StatusBadRequest, errors.New("name cannot be empty"))
		return
	}

	req := mapper.CompanyUpdateReqToUC(id, dtoReq)

	err = h.update.Execute(ctx, req)
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
// @Failure 404 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
// @Router /companies/{id} [delete]
func (h *CompanyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := middleware.UUIDFromContext(ctx)
	if err != nil {
		h.handleErr(w, err)
		return
	}

	err = h.delete.Execute(ctx, id)
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

	case errors.Is(err, domain_errors.ErrCompanyAlreadyExists):
		responds.RespondError(w, http.StatusConflict, err)

	default:
		responds.RespondError(w, http.StatusInternalServerError, err)
	}
}
