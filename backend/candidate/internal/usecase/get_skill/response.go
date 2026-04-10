package get_skill

import (
	"github.com/google/uuid"
)

type GetByIdResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ListResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
