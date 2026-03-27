package list

import (
	"time"

	"github.com/google/uuid"
)

type CompanySummary struct {
	ID               uuid.UUID `db:"id"`
	Name             string    `db:"name"`
	OpenVacanciesCnt int       `db:"open_vacancies_count"`
	LogoKey          *string   `db:"logo_key"`
	CreatedAt        time.Time `db:"created_at"`
}

type Response struct {
	Companies  []CompanySummary
	NextCursor *string
}
