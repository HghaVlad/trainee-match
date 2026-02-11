package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/auth"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/helpers"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/create_candidate"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/get_candidate"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/get_candidate_by_user_id"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/update_candidate"
	"net/http"
	"time"
)

type Candidate struct {
	getById     *get_candidate.UseCase
	create      *create_candidate.UseCase
	update      *update_candidate.UseCase
	getByUserId *get_candidate_by_user_id.UseCase
}

func NewCandidate(getById *get_candidate.UseCase, create *create_candidate.UseCase, update *update_candidate.UseCase, getByUserId *get_candidate_by_user_id.UseCase) *Candidate {
	return &Candidate{
		getById:     getById,
		create:      create,
		update:      update,
		getByUserId: getByUserId,
	}
}

// GetMe godoc
// @Summary Get my candidate profile
// @Tags candidate
// @Accept json
// @Produce json
// @Success 200 {object} dto.CandidateResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /candidate/me [get]
func (c *Candidate) GetMe(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.FromContext(r.Context())
	if !ok {
		helpers.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	candidate, err := c.getByUserId.Execute(r.Context(), user.Id)
	if errors.Is(err, domain.ErrCandidateNotFound) {
		helpers.RespondError(w, http.StatusNotFound, "candidate not found")
		return
	} else if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.RespondJSON(w, http.StatusOK, dto.CandidateResponse{
		ID:       candidate.ID,
		UserID:   candidate.UserID,
		Phone:    candidate.Phone,
		Telegram: candidate.Telegram,
		City:     candidate.City,
		Birthday: candidate.Birthday.Format("02.01.2006"),
	})
}

// CreateCandidate godoc
// @Summary Create candidate profile
// @Description Creates a new candidate profile associated with the authenticated user
// @Tags candidate
// @Accept json
// @Produce json
// @Param input body dto.CandidateCreateRequest true "Candidate creation data"
// @Success 201 {object} dto.CandidateResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /candidate/ [post]
func (c *Candidate) CreateCandidate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req dto.CandidateCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body %e", err))
		return
	}
	if err := req.Validate(); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body %e", err))
		return
	}

	user, ok := auth.FromContext(r.Context())
	if !ok {
		helpers.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	_, err := c.getByUserId.Execute(r.Context(), user.Id)
	if err == nil {
		helpers.RespondError(w, http.StatusConflict, "candidate already exists")
		return
	}

	candidate, err := c.create.Execute(r.Context(), &create_candidate.Request{
		UserID:   user.Id,
		Phone:    req.Phone,
		Telegram: req.Telegram,
		City:     req.City,
		Birthday: time.Time(req.Birthday),
	})
	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.CandidateResponse{
		ID:       candidate,
		UserID:   user.Id,
		Phone:    req.Phone,
		Telegram: req.Telegram,
		City:     req.City,
		Birthday: time.Time(req.Birthday).Format("02.01.2006"),
	}

	helpers.RespondJSON(w, http.StatusCreated, response)
}

// UpdateCandidate godoc
// @Summary Update candidate profile
// @Tags candidate
// @Accept json
// @Produce json
// @Param input body dto.CandidateUpdateRequest true "Candidate creation data"
// @Success 200 {object} dto.CandidateResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /candidate/ [patch]
func (c *Candidate) UpdateCandidate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req dto.CandidateUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body %e", err))
		return
	}
	if err := req.Validate(); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body %e", err))
		return
	}

	user, ok := auth.FromContext(r.Context())
	if !ok {
		helpers.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	candidate, err := c.getByUserId.Execute(r.Context(), user.Id)
	if errors.Is(err, domain.ErrCandidateNotFound) {
		helpers.RespondError(w, http.StatusNotFound, "candidate not found")
	} else if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var birthday *time.Time
	if req.Birthday != nil {
		t := time.Time(*req.Birthday)
		birthday = &t
	}

	updatedCandidate, err := c.update.Execute(r.Context(), &update_candidate.Request{
		ID:       candidate.ID,
		UserID:   &user.Id,
		Phone:    req.Phone,
		Telegram: req.Telegram,
		City:     req.City,
		Birthday: birthday,
	})

	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	helpers.RespondJSON(w, http.StatusOK, dto.CandidateResponse{
		ID:       updatedCandidate.ID,
		UserID:   updatedCandidate.UserID,
		Phone:    updatedCandidate.Phone,
		Telegram: updatedCandidate.Telegram,
		City:     updatedCandidate.City,
		Birthday: updatedCandidate.Birthday.Format("02.01.2006"),
	})

}
