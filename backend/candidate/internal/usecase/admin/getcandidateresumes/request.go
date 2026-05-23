package getcandidateresumes

import "github.com/google/uuid"

type Request struct {
	CandidateID uuid.UUID
	Page        int
	Size        int
}
