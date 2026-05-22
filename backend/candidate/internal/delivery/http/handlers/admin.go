package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/helpers"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/admin/addskill"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/admin/archiveresume"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/admin/deleteskill"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/admin/getcandidate"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/admin/getcandidateresumes"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/admin/getcandidates"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/admin/getresume"
)

type Admin struct {
	getCandidatesUC       *getcandidates.UseCase
	getCandidateUC        *getcandidate.UseCase
	getCandidateResumesUC *getcandidateresumes.UseCase
	getResumeUC           *getresume.UseCase
	archiveResumeUC       *archiveresume.UseCase
	addSkillUC            *addskill.UseCase
	deleteSkillUC         *deleteskill.UseCase
}

func NewAdmin(
	getCandidatesUC *getcandidates.UseCase,
	getCandidateUC *getcandidate.UseCase,
	getCandidateResumesUC *getcandidateresumes.UseCase,
	getResumeUC *getresume.UseCase,
	archiveResumeUC *archiveresume.UseCase,
	addSkillUC *addskill.UseCase,
	deleteSkillUC *deleteskill.UseCase,
) *Admin {
	return &Admin{
		getCandidatesUC:       getCandidatesUC,
		getCandidateUC:        getCandidateUC,
		getCandidateResumesUC: getCandidateResumesUC,
		getResumeUC:           getResumeUC,
		archiveResumeUC:       archiveResumeUC,
		addSkillUC:            addSkillUC,
		deleteSkillUC:         deleteSkillUC,
	}
}

// GetCandidates godoc
// @Summary List candidates (admin)
// @Tags admin
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {array} dto.CandidateResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/candidates [get]
func (a *Admin) GetCandidates(w http.ResponseWriter, r *http.Request) {
	page, size, err := helpers.ParsePageSize(r)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	candidates, err := a.getCandidatesUC.Execute(r.Context(), getcandidates.Request{Page: page, Size: size})
	if err != nil {
		helpers.RespondErrorSmart(w, err)
		return
	}

	resp := make([]dto.CandidateResponse, len(candidates))
	for i, c := range candidates {
		resp[i] = dto.CandidateResponse{
			ID:       c.ID,
			UserID:   c.UserID,
			FullName: c.FullName,
			Phone:    c.Phone,
			Telegram: c.Telegram,
			City:     c.City,
			Birthday: dto.TimeToDate(c.Birthday),
		}
	}

	helpers.RespondJSON(w, http.StatusOK, resp)
}

// GetCandidate godoc
// @Summary Get candidate by ID (admin)
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Candidate ID"
// @Success 200 {object} dto.CandidateResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/candidates/{id} [get]
func (a *Admin) GetCandidate(w http.ResponseWriter, r *http.Request) {
	candidateID, err := helpers.ParseUUIDParam(r, "id")
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	candidate, err := a.getCandidateUC.Execute(r.Context(), getcandidate.Request{CandidateID: candidateID})
	if err != nil {
		helpers.RespondErrorSmart(w, err)
		return
	}

	resp := dto.CandidateResponse{
		ID:       candidate.ID,
		UserID:   candidate.UserID,
		FullName: candidate.FullName,
		Phone:    candidate.Phone,
		Telegram: candidate.Telegram,
		City:     candidate.City,
		Birthday: dto.TimeToDate(candidate.Birthday),
	}

	helpers.RespondJSON(w, http.StatusOK, resp)
}

// GetCandidateResumes godoc
// @Summary List resumes for candidate (admin)
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Candidate ID"
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Success 200 {array} dto.ShortResumeResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/candidates/{id}/resumes [get]
func (a *Admin) GetCandidateResumes(w http.ResponseWriter, r *http.Request) {
	candidateID, err := helpers.ParseUUIDParam(r, "id")
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	page, size, err := helpers.ParsePageSize(r)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	resumes, err := a.getCandidateResumesUC.Execute(r.Context(), getcandidateresumes.Request{
		CandidateID: candidateID,
		Page:        page,
		Size:        size,
	})
	if err != nil {
		helpers.RespondErrorSmart(w, err)
		return
	}

	resp := make([]dto.ShortResumeResponse, len(resumes))
	for i, c := range resumes {
		resp[i] = dto.ShortResumeResponse{
			ID:               c.ID,
			CandidateId:      c.CandidateID,
			Name:             c.Name,
			Status:           c.Status,
			ModerationStatus: c.ModerationStatus,
		}
	}

	helpers.RespondJSON(w, http.StatusOK, resp)
}

// GetResume godoc
// @Summary Get resume by ID (admin)
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Resume ID"
// @Success 200 {object} dto.ResumeResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/resumes/{id} [get]
func (a *Admin) GetResume(w http.ResponseWriter, r *http.Request) {
	resumeID, err := helpers.ParseUUIDParam(r, "id")
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	resume, err := a.getResumeUC.Execute(r.Context(), getresume.Request{ResumeID: resumeID})
	if err != nil {
		helpers.RespondErrorSmart(w, err)
		return
	}
	resp := dto.AdminUseCaseResponseToDtoResumeResponse(resume)

	helpers.RespondJSON(w, http.StatusOK, resp)
}

// ArchiveResume godoc
// @Summary Archive resume (admin)
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Resume ID"
// @Success 204 {string} string "ok"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/resumes/{id}/archive [post]
func (a *Admin) ArchiveResume(w http.ResponseWriter, r *http.Request) {
	resumeID, err := helpers.ParseUUIDParam(r, "id")
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = a.archiveResumeUC.Execute(r.Context(), archiveresume.Request{ResumeID: resumeID})
	if err != nil {
		helpers.RespondErrorSmart(w, err)
		return
	}

	helpers.RespondJSON(w, http.StatusOK, struct{ Message string }{Message: "ok"})
}

// AddSkill godoc
// @Summary Create skill (admin)
// @Tags admin
// @Accept json
// @Produce json
// @Param input body dto.SkillRequest true "Skill creation data"
// @Success 201 {object} dto.SkillResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/skills [post]
func (a *Admin) AddSkill(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req dto.SkillRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := a.addSkillUC.Execute(r.Context(), addskill.Request{Name: req.Name})
	if err != nil {
		helpers.RespondErrorSmart(w, err)
		return
	}

	helpers.RespondJSON(w, http.StatusCreated, dto.SkillResponse{ID: resp.ID, Name: resp.Name})
}

// DeleteSkill godoc
// @Summary Delete skill (admin)
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Skill ID"
// @Success 204 {string} string "ok"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/skills/{id} [delete]
func (a *Admin) DeleteSkill(w http.ResponseWriter, r *http.Request) {
	skillID, err := helpers.ParseUUIDParam(r, "id")
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = a.deleteSkillUC.Execute(r.Context(), deleteskill.Request{SkillID: skillID})
	if err != nil {
		helpers.RespondErrorSmart(w, err)
		return
	}

	helpers.RespondJSON(w, http.StatusOK, struct{ Message string }{Message: "ok"})
}
