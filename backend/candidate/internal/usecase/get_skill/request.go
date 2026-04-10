package get_skill

import (
	"github.com/google/uuid"
)

type GetByIdRequest struct {
	ID uuid.UUID `json:"id"`
}

type ListRequest struct {
	// Empty struct for now, can be expanded with filters/pagination if needed
}
