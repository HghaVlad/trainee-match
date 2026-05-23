package getcandidate

import "github.com/google/uuid"

type Request struct {
	CandidateID uuid.UUID `json:"candidate_id"`
}
