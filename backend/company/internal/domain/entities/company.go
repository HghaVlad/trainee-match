package domain

import (
	"time"

	"github.com/google/uuid"
)

type Company struct {
	ID               uuid.UUID `db:"id"`
	Name             string    `db:"name"`
	OpenVacanciesCnt int       `db:"open_vacancies_count"`
	Description      *string   `db:"description"`
	Website          *string   `db:"website"`
	LogoKey          *string   `db:"logo_key"`
	OwnerID          uuid.UUID `db:"owner_id"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
