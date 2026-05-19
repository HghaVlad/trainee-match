package memget

import (
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
)

type Response struct {
	ID               uuid.UUID
	Name             string
	OpenVacanciesCnt int
	Description      *string
	Website          *string
	LogoURL          *string
	ModStatus        company.ModerationStatus
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
