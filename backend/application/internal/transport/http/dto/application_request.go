package dto

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrBadRequest = errors.New("bad request")
)

type CreateApplicationRequest struct {
	ResumeID  uuid.UUID `json:"resume_id"`
	VacancyID uuid.UUID `json:"vacancy_id"`
}
