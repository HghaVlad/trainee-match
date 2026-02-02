package domain

import (
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
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

const (
	MaxCompanyNameLen        = 80
	MaxCompanyDescriptionLen = 5000
)

func (c *Company) Validate() error {
	if len(c.Name) == 0 || len(c.Name) > MaxCompanyNameLen {
		return domain_errors.ErrCompanyInvalidNameLen
	}

	if c.Description != nil && len(*c.Description) > MaxCompanyDescriptionLen {
		return domain_errors.ErrCompanyInvalidDescriptionLen
	}

	return nil
}
