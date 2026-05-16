package companyupdated

import (
	"context"

	"github.com/google/uuid"
)

type VacancyRepo interface {
	UpdateCompanyName(ctx context.Context, companyID uuid.UUID, companyName string) error
}
