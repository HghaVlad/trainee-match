package handlers

import (
	"encoding/json"
	"errors"
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
// @Product json
// @Success 200 {object} dto.CandidateResponse
func (c *Candidate) GetMe(w http.ResponseWriter, r *http.Request) {
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

	helpers.RespondJSON(w, http.StatusCreated, dto.CandidateResponse{
		ID:       candidate.ID,
		UserID:   candidate.UserID,
		Phone:    candidate.Phone,
		Telegram: candidate.Telegram,
		City:     candidate.City,
		Birthday: candidate.Birthday,
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
// @Failure 400 {object} responds.ErrorResponse
// @Failure 409 {object} responds.ErrorResponse
// @Failure 500 {object} responds.ErrorResponse
func (c *Candidate) CreateCandidate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req dto.CandidateCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	birthday, err := time.Parse("02.01.2006", req.Birthday)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "invalid birthday format")
		return
	}

	user, ok := auth.FromContext(r.Context())
	if !ok {
		helpers.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	_, err = c.getByUserId.Execute(r.Context(), user.Id)
	if err == nil {
		helpers.RespondError(w, http.StatusConflict, "candidate already exists")
		return
	}

	candidate, err := c.create.Execute(r.Context(), &create_candidate.Request{
		UserID:   user.Id,
		Phone:    req.Phone,
		Telegram: req.Telegram,
		City:     req.City,
		Birthday: birthday,
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
		Birthday: birthday,
	}

	helpers.RespondJSON(w, http.StatusCreated, response)
}
