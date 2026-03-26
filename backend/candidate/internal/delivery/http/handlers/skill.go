package handlers

import (
	"errors"
	"net/http"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/delivery/http/helpers"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/get_skill"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Skill struct {
	getSkillUC *get_skill.UseCase
}

func NewSkill(getSkillUC *get_skill.UseCase) *Skill {
	return &Skill{
		getSkillUC: getSkillUC,
	}
}

// GetSkill godoc
// @Summary Get skill by ID
// @Tags skill
// @Accept json
// @Produce json
// @Param id path string true "Skill ID"
// @Success 200 {object} dto.SkillResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /skill/{id} [get]
func (s *Skill) GetSkill(w http.ResponseWriter, r *http.Request) {
	// Extract skill ID from URL path
	skillId := chi.URLParam(r, "id")
	if skillId == "" {
		helpers.RespondError(w, http.StatusBadRequest, "skill ID is required")
		return
	}

	parsedId, err := uuid.Parse(skillId)
	if err != nil {
		helpers.RespondError(w, http.StatusBadRequest, "invalid skill ID format")
		return
	}

	req := get_skill.GetByIdRequest{ID: parsedId}
	skill, err := s.getSkillUC.Execute(r.Context(), req)
	if errors.Is(err, domain.ErrSkillNotFound) {
		helpers.RespondError(w, http.StatusNotFound, "skill not found")
		return
	} else if errors.Is(err, domain.ErrInvalidSkillName) {
		helpers.RespondError(w, http.StatusBadRequest, err.Error())
		return
	} else if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.SkillResponse{
		ID:   skill.ID,
		Name: skill.Name,
	}

	helpers.RespondJSON(w, http.StatusOK, response)
}

// ListSkills godoc
// @Summary List all skills
// @Tags skill
// @Accept json
// @Produce json
// @Success 200 {array} dto.SkillResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /skill/list [get]
func (s *Skill) ListSkills(w http.ResponseWriter, r *http.Request) {
	req := get_skill.ListRequest{}
	skills, err := s.getSkillUC.ExecuteList(r.Context(), req)
	if err != nil {
		helpers.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var responses []dto.SkillResponse
	for _, skill := range skills {
		responses = append(responses, dto.SkillResponse{
			ID:   skill.ID,
			Name: skill.Name,
		})
	}

	helpers.RespondJSON(w, http.StatusOK, responses)
}
