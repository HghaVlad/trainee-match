package get_resume

import (
	"github.com/google/uuid"
)

type GetByIdRequest struct {
	ID uuid.UUID `json:"id"`
}

type GetByCandidateIdRequest struct {
	CandidateId uuid.UUID `json:"candidate_id"`
}