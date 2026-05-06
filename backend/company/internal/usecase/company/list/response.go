package list

import (
	"time"

	"github.com/google/uuid"
)

type CompanySummary struct {
	ID               uuid.UUID
	Name             string
	OpenVacanciesCnt int
	LogoKey          *string
	CreatedAt        time.Time
}

type Response struct {
	Companies  []CompanySummary
	NextCursor *string
}
