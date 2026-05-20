package create

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

type VacancyRepo interface {
	Create(ctx context.Context, vacancy *vacancy.Vacancy) error
}

type CompMemberRepo interface {
	Get(ctx context.Context, userID, companyID uuid.UUID) (*member.CompanyMember, error)
}

type CompanyRepo interface {
	GetByMember(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*company.Company, error)
}

type SearchRepo interface {
	Index(ctx context.Context, vac vacancy.Vacancy, compName string) error
}
