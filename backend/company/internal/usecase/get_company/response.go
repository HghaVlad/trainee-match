package get_company

import (
	"time"

	"github.com/google/uuid"
)

type Response struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Website     *string
	LogoURL     *string // IDK about that
	OwnerID     uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
