package company

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
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

const (
	MaxCompanyNameLen        = 80
	MaxCompanyDescriptionLen = 5000
)

func (c *Company) Validate() error {
	if len(c.Name) == 0 || len([]rune(c.Name)) > MaxCompanyNameLen {
		return ErrCompanyInvalidNameLen
	}

	if c.Description != nil && len([]rune(*c.Description)) > MaxCompanyDescriptionLen {
		return ErrCompanyInvalidDescriptionLen
	}

	return nil
}
