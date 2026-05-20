package companydeleted

import (
	"context"

	"github.com/google/uuid"
)

type CompanyMemberRepo interface {
	DeleteByCompanyID(ctx context.Context, companyID uuid.UUID) error
}

type VacancyRepo interface {
	DeleteByCompanyID(ctx context.Context, companyID uuid.UUID) error
}
