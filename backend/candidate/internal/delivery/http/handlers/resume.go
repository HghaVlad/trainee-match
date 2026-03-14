package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/auth"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/helpers"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/create_resume"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/get_resume"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/update_resume"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Resume struct {
	createResumeUC *create_resume.UseCase
	getResumeUC    *get_resume.UseCase
	updateResumeUC *update_resume.UseCase
}

func NewResume(createResumeUC *create_resume.UseCase, getResumeUC *get_resume.UseCase,
	updateResumeUC *update_resume.UseCase) *Resume {
	return &Resume{
		createResumeUC: createResumeUC,
		getResumeUC:    getResumeUC,
		updateResumeUC: updateResumeUC,
	}
}

// CreateResume godoc
// @Summary Create a new resume
// @Tags resume
// @Accept json
// @Produce json
// @Param input body dto.CreateResumeRequest true "Resume creation data"
// @Success 201 {object} dto.ResumeResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /resume/ [post]
func (res *Resume) CreateResume(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req dto.CreateResumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body %v", err))
		return
	}

	if err := req.Validate(); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("validation error: %s", err.Error()))
		return
	}

	user, ok := auth.FromContext(r.Context())
	if !ok {
		helpers.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	useCaseRequest := req.ToUseCaseRequest()
	useCaseRequest.UserId = user.Id

	resp, err := res.createResumeUC.Execute(r.Context(), useCaseRequest)
	if err != nil {
		helpers.RespondErrorSmart(w, err)
		return
	}

	response := dto.ResumeResponse{
		ID:          resp.ID,
		CandidateID: resp.CandidateID,
		Name:        req.Name,
		Status:      req.Status,
		Data:        req.Data,
	}

	helpers.RespondJSON(w, http.StatusCreated, response)
}

// GetResume godoc
// @Summary Get resume by ID
// @Tags resume
// @Accept json
// @Produce json
// @Param id path string true "Resume ID"
// @Success 200 {object} dto.ResumeResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /resume/{id} [get]
func (res *Resume) GetResume(w http.ResponseWriter, r *http.Request) {
	resumeId := chi.URLParam(r, "id")
	if resumeId == "" {
		helpers.RespondError(w, http.StatusBadRequest, "resume ID is required")
		return
	}

	parsedId, err := uuid.Parse(resumeId)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "invalid resume ID format")
		return
	}

	user, ok := auth.FromContext(r.Context())
	if !ok {
		helpers.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	useCaseResp, err := res.getResumeUC.GetById(r.Context(), parsedId, user.Id)
	if err != nil {
		helpers.RespondErrorSmart(w, err)
		return
	}

	response := dto.UseCaseResponseToDtoResumeResponse(*useCaseResp)

	helpers.RespondJSON(w, http.StatusOK, response)
}

// UpdateResume godoc
// @Summary Update a resume
// @Tags resume
// @Accept json
// @Produce json
// @Param id path string true "Resume ID"
// @Param input body dto.UpdateResumeRequest true "Resume update data"
// @Success 200 {object} dto.ResumeResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /resume/{id} [patch]
func (res *Resume) UpdateResume(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	resumeId := chi.URLParam(r, "id")
	if resumeId == "" {
		helpers.RespondError(w, http.StatusBadRequest, "resume ID is required")
		return
	}

	parsedId, err := uuid.Parse(resumeId)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "invalid resume ID format")
		return
	}

	var req dto.UpdateResumeRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body %e", err))
		return
	}

	// Validate the request
	if err = req.Validate(); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("validation error: %s", err.Error()))
		return
	}

	user, ok := auth.FromContext(r.Context())
	if !ok {
		helpers.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Prepare the use case request
	useCaseReq := req.ToUseCaseRequest()
	useCaseReq.ID = parsedId
	useCaseReq.UserId = user.Id

	err = res.updateResumeUC.Execute(r.Context(), useCaseReq)
	if err != nil {
		helpers.RespondErrorSmart(w, err)
		return
	}

	// Return the updated resume
	updatedResumeResp, err := res.getResumeUC.GetById(r.Context(), parsedId, user.Id)
	if err != nil {
		helpers.RespondErrorSmart(w, err)
		return
	}
	response := dto.UseCaseResponseToDtoResumeResponse(*updatedResumeResp)

	helpers.RespondJSON(w, http.StatusOK, response)
}

// ListResumes godoc
// @Summary List all resumes for the authenticated candidate
// @Tags resume
// @Accept json
// @Produce json
// @Success 200 {object} []dto.ShortResumeResponse
// @Failure 401 {object} dto.ErrorResponse "unauthorized"
// @Failure 404 {object} dto.ErrorResponse "candidate not found"
// @Failure 500 {object} dto.ErrorResponse
// @Router /resume [get]
func (res *Resume) ListResumes(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.FromContext(r.Context())
	if !ok {
		helpers.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	useCaseResp, err := res.getResumeUC.GetByCandidateId(r.Context(), user.Id)
	if err != nil {
		helpers.RespondErrorSmart(w, err)
		return
	}

	responses := make([]dto.ShortResumeResponse, len(useCaseResp))
	for i, resume := range useCaseResp {
		responses[i] = dto.ShortResumeResponse{
			ID:          resume.ID,
			CandidateId: resume.CandidateId,
			Name:        resume.Name,
			Status:      resume.Status,
		}
	}

	helpers.RespondJSON(w, http.StatusOK, responses)
}
