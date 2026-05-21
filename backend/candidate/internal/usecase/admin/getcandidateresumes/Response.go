package getcandidateresumes

import "github.com/google/uuid"

type ShortResponse struct {
	ID               uuid.UUID `json:"id"`
	CandidateID      uuid.UUID `json:"candidate_id"`
	Name             string    `json:"name"`
	Status           string    `json:"status"`
	ModerationStatus string    `json:"moderation_status"`
}
