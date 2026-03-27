package update

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
)

type Request struct {
	ID          uuid.UUID
	Name        *string
	Description *string
	Website     *string
}

func (c *Request) Validate() error {
	if c.Name != nil && (len(*c.Name) == 0 || len(*c.Name) > company.MaxCompanyNameLen) {
		return company.ErrCompanyInvalidNameLen
	}

	if c.Description != nil && len(*c.Description) > company.MaxCompanyDescriptionLen {
		return company.ErrCompanyInvalidDescriptionLen
	}

	return nil
}
