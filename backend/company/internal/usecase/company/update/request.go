package update_company

import (
	"github.com/google/uuid"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	domain_errors "github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
)

type Request struct {
	ID          uuid.UUID
	Name        *string
	Description *string
	Website     *string
}

func (c *Request) Validate() error {
	if c.Name != nil && (len(*c.Name) == 0 || len(*c.Name) > domain.MaxCompanyNameLen) {
		return domain_errors.ErrCompanyInvalidNameLen
	}

	if c.Description != nil && len(*c.Description) > domain.MaxCompanyDescriptionLen {
		return domain_errors.ErrCompanyInvalidDescriptionLen
	}

	return nil
}
