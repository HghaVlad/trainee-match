package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/auth"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/helpers"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/create_resume"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/get_resume"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/update_resume"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Resume struct {
	createResumeUC *create_resume.UseCase
	getResumeUC    *get_resume.UseCase
	updateResumeUC *update_resume.UseCase
}

func NewResume(createResumeUC *create_resume.UseCase, getResumeUC *get_resume.UseCase, updateResumeUC *update_resume.UseCase) *Resume {
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
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body %e", err))
		return
	}

	// Validate the request
	if err := req.Validate(); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("validation error: %s", err.Error()))
		return
	}

	user, ok := auth.FromContext(r.Context())
	if !ok {
		helpers.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	useCaseReq := create_resume.Request{
		CandidateId: user.Id,
		Name:        req.Name,
		Status:      req.Status,
		Data: create_resume.ResumeData{
			LastName:        req.Data.LastName,
			FirstName:       req.Data.FirstName,
			MiddleName:      req.Data.MiddleName,
			DateOfBirth:     time.Time(req.Data.DateOfBirth).Format("02.01.2006"),
			Email:           req.Data.Email,
			Phone:           req.Data.Phone,
			City:            req.Data.City,
			Citizenship:     req.Data.Citizenship,
			Education:       make([]create_resume.Education, len(req.Data.Education)),
			WorkExperiences: make([]create_resume.WorkExperience, len(req.Data.WorkExperiences)),
			SkillsList:      req.Data.SkillsList,
			AdditionalInfo:  req.Data.AdditionalInfo,
			PortfolioLink:   req.Data.PortfolioLink,
			DesiredFormat:   req.Data.DesiredFormat,
			EnglishLevel:    req.Data.EnglishLevel,
		},
	}

	for i, edu := range req.Data.Education {
		useCaseReq.Data.Education[i] = create_resume.Education{
			Level:          edu.Level,
			University:     edu.University,
			Faculty:        edu.Faculty,
			Specialization: edu.Specialization,
			StartYear:      edu.StartYear,
			EndYear:        edu.EndYear,
			Format:         edu.Format,
		}
	}

	for i, exp := range req.Data.WorkExperiences {
		useCaseReq.Data.WorkExperiences[i] = create_resume.WorkExperience{
			Position:         exp.Position,
			Company:          exp.Company,
			Period:           exp.Period,
			Responsibilities: exp.Responsibilities,
		}
	}

	resp, err := res.createResumeUC.Execute(r.Context(), useCaseReq)
	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.ResumeResponse{
		ID:          resp.ID,
		CandidateId: user.Id,
		Name:        req.Name,
		Status:      req.Status,
		Data:        req.Data,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
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

	req := get_resume.GetByIdRequest{ID: parsedId}
	useCaseResp, err := res.getResumeUC.GetById(r.Context(), req)
	if errors.Is(err, domain.ErrResumeNotFound) {
		helpers.RespondError(w, http.StatusNotFound, "resume not found")
		return
	} else if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check if the user has access to this resume (must be the owner)
	user, ok := auth.FromContext(r.Context())
	if !ok {
		helpers.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if useCaseResp.CandidateId != user.Id {
		helpers.RespondError(w, http.StatusForbidden, "forbidden")
		return
	}

	// Convert use case response to DTO
	educationDTO := make([]dto.Education, len(useCaseResp.Data.Education))
	for i, edu := range useCaseResp.Data.Education {
		educationDTO[i] = dto.Education{
			Level:          edu.Level,
			University:     edu.University,
			Faculty:        edu.Faculty,
			Specialization: edu.Specialization,
			StartYear:      edu.StartYear,
			EndYear:        edu.EndYear,
			Format:         edu.Format,
		}
	}

	workExpDTO := make([]dto.WorkExperience, len(useCaseResp.Data.WorkExperiences))
	for i, exp := range useCaseResp.Data.WorkExperiences {
		workExpDTO[i] = dto.WorkExperience{
			Position:         exp.Position,
			Company:          exp.Company,
			Period:           exp.Period,
			Responsibilities: exp.Responsibilities,
		}
	}

	// Parse date string back to DTO Date type
	dtoDateOfBirth := dto.Date{}
	dateStr := fmt.Sprintf("\"%s\"", useCaseResp.Data.DateOfBirth)
	err = dtoDateOfBirth.UnmarshalJSON([]byte(dateStr))
	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, "failed to parse date of birth")
		return
	}

	dtoData := dto.ResumeData{
		LastName:        useCaseResp.Data.LastName,
		FirstName:       useCaseResp.Data.FirstName,
		MiddleName:      useCaseResp.Data.MiddleName,
		DateOfBirth:     dtoDateOfBirth,
		Email:           useCaseResp.Data.Email,
		Phone:           useCaseResp.Data.Phone,
		City:            useCaseResp.Data.City,
		Citizenship:     useCaseResp.Data.Citizenship,
		Education:       educationDTO,
		WorkExperiences: workExpDTO,
		SkillsList:      useCaseResp.Data.SkillsList,
		AdditionalInfo:  useCaseResp.Data.AdditionalInfo,
		PortfolioLink:   useCaseResp.Data.PortfolioLink,
		DesiredFormat:   useCaseResp.Data.DesiredFormat,
		EnglishLevel:    useCaseResp.Data.EnglishLevel,
	}

	response := dto.ResumeResponse{
		ID:          useCaseResp.ID,
		CandidateId: useCaseResp.CandidateId,
		Name:        useCaseResp.Name,
		Status:      useCaseResp.Status,
		Data:        dtoData,
		CreatedAt:   useCaseResp.CreatedAt,
		UpdatedAt:   useCaseResp.UpdatedAt,
	}

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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body %e", err))
		return
	}

	// Validate the request
	if err := req.Validate(); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("validation error: %s", err.Error()))
		return
	}

	// Check if the user has access to this resume (must be the owner)
	getReq := get_resume.GetByIdRequest{ID: parsedId}
	resume, err := res.getResumeUC.GetById(r.Context(), getReq)
	if errors.Is(err, domain.ErrResumeNotFound) {
		helpers.RespondError(w, http.StatusNotFound, "resume not found")
		return
	} else if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	user, ok := auth.FromContext(r.Context())
	if !ok {
		helpers.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if resume.CandidateId != user.Id {
		helpers.RespondError(w, http.StatusForbidden, "forbidden")
		return
	}

	// Prepare the use case request
	useCaseReq := update_resume.Request{
		ID:     parsedId,
		Name:   req.Name,
		Status: req.Status,
	}

	if req.Data != nil {
		// Parse DateOfBirth from DTO.Date to string
		dateOfBirthBytes, err := req.Data.DateOfBirth.MarshalJSON()
		if err != nil {
			helpers.RespondError(w, http.StatusInternalServerError, "failed to marshal date of birth")
			return
		}
		var dateStr string
		err = json.Unmarshal(dateOfBirthBytes, &dateStr)
		if err != nil {
			helpers.RespondError(w, http.StatusInternalServerError, "failed to unmarshal date of birth")
			return
		}

		useCaseReq.Data = &update_resume.ResumeData{
			LastName:        req.Data.LastName,
			FirstName:       req.Data.FirstName,
			MiddleName:      req.Data.MiddleName,
			DateOfBirth:     dateStr, // Convert Date to string
			Email:           req.Data.Email,
			Phone:           req.Data.Phone,
			City:            req.Data.City,
			Citizenship:     req.Data.Citizenship,
			Education:       make([]update_resume.Education, len(req.Data.Education)),
			WorkExperiences: make([]update_resume.WorkExperience, len(req.Data.WorkExperiences)),
			SkillsList:      req.Data.SkillsList,
			AdditionalInfo:  req.Data.AdditionalInfo,
			PortfolioLink:   req.Data.PortfolioLink,
			DesiredFormat:   req.Data.DesiredFormat,
			EnglishLevel:    req.Data.EnglishLevel,
		}

		for i, edu := range req.Data.Education {
			useCaseReq.Data.Education[i] = update_resume.Education{
				Level:          edu.Level,
				University:     edu.University,
				Faculty:        edu.Faculty,
				Specialization: edu.Specialization,
				StartYear:      edu.StartYear,
				EndYear:        edu.EndYear,
				Format:         edu.Format,
			}
		}

		for i, exp := range req.Data.WorkExperiences {
			useCaseReq.Data.WorkExperiences[i] = update_resume.WorkExperience{
				Position:         exp.Position,
				Company:          exp.Company,
				Period:           exp.Period,
				Responsibilities: exp.Responsibilities,
			}
		}
	}

	_, err = res.updateResumeUC.Execute(r.Context(), useCaseReq)
	if err != nil {
		if errors.Is(err, domain.ErrResumeNotFound) {
			helpers.RespondError(w, http.StatusNotFound, "resume not found")
			return
		} else if errors.Is(err, domain.ErrSkillNotFound) {
			helpers.RespondError(w, http.StatusBadRequest, "one or more skills not found")
			return
		}
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return the updated resume
	updatedResumeReq := get_resume.GetByIdRequest{ID: parsedId}
	updatedResumeResp, err := res.getResumeUC.GetById(r.Context(), updatedResumeReq)
	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert use case response to DTO
	educationDTO := make([]dto.Education, len(updatedResumeResp.Data.Education))
	for i, edu := range updatedResumeResp.Data.Education {
		educationDTO[i] = dto.Education{
			Level:          edu.Level,
			University:     edu.University,
			Faculty:        edu.Faculty,
			Specialization: edu.Specialization,
			StartYear:      edu.StartYear,
			EndYear:        edu.EndYear,
			Format:         edu.Format,
		}
	}

	workExpDTO := make([]dto.WorkExperience, len(updatedResumeResp.Data.WorkExperiences))
	for i, exp := range updatedResumeResp.Data.WorkExperiences {
		workExpDTO[i] = dto.WorkExperience{
			Position:         exp.Position,
			Company:          exp.Company,
			Period:           exp.Period,
			Responsibilities: exp.Responsibilities,
		}
	}

	// Parse date string back to DTO Date type
	dtoDateOfBirth := dto.Date{}
	dateStr := fmt.Sprintf("\"%s\"", updatedResumeResp.Data.DateOfBirth)
	err = dtoDateOfBirth.UnmarshalJSON([]byte(dateStr))
	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, "failed to parse date of birth")
		return
	}

	dtoData := dto.ResumeData{
		LastName:        updatedResumeResp.Data.LastName,
		FirstName:       updatedResumeResp.Data.FirstName,
		MiddleName:      updatedResumeResp.Data.MiddleName,
		DateOfBirth:     dtoDateOfBirth,
		Email:           updatedResumeResp.Data.Email,
		Phone:           updatedResumeResp.Data.Phone,
		City:            updatedResumeResp.Data.City,
		Citizenship:     updatedResumeResp.Data.Citizenship,
		Education:       educationDTO,
		WorkExperiences: workExpDTO,
		SkillsList:      updatedResumeResp.Data.SkillsList,
		AdditionalInfo:  updatedResumeResp.Data.AdditionalInfo,
		PortfolioLink:   updatedResumeResp.Data.PortfolioLink,
		DesiredFormat:   updatedResumeResp.Data.DesiredFormat,
		EnglishLevel:    updatedResumeResp.Data.EnglishLevel,
	}

	response := dto.ResumeResponse{
		ID:          updatedResumeResp.ID,
		CandidateId: updatedResumeResp.CandidateId,
		Name:        updatedResumeResp.Name,
		Status:      updatedResumeResp.Status,
		Data:        dtoData,
		CreatedAt:   updatedResumeResp.CreatedAt,
		UpdatedAt:   updatedResumeResp.UpdatedAt,
	}

	helpers.RespondJSON(w, http.StatusOK, response)
}
