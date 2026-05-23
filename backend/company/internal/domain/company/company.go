package company

import (
	"time"

	"github.com/google/uuid"
)

type Company struct {
	ID               uuid.UUID
	Name             string
	OpenVacanciesCnt int
	Description      *string
	Website          *string
	LogoKey          *string
	ModerationStatus ModerationStatus
	CreatedAt        time.Time
	UpdatedAt        time.Time
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

type ModerationStatus string

const (
	ModerationStatusOK     ModerationStatus = "ok"
	ModerationStatusHidden ModerationStatus = "hidden"
)

func (ms ModerationStatus) IsValid() bool {
	switch ms {
	case ModerationStatusOK,
		ModerationStatusHidden:
		return true
	}

	return false
}
