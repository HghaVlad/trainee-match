package get_skill

import (
	"github.com/google/uuid"
)

type GetByIdRequest struct {
	ID uuid.UUID `json:"id"`
}

type ListRequest struct {
	Page int
	Size int
}
