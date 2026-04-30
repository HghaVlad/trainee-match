package publish

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

type PublishedEventView struct {
	ID                  uuid.UUID
	Title               string
	CompanyID           uuid.UUID
	CompanyName         string
	Status              vacancy.Status
	WasAlreadyPublished bool
}
