package publish

import (
	"context"

	"github.com/google/uuid"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

type VacancyRepo interface {
	GetByID(ctx context.Context, vacancyID uuid.UUID, companyID uuid.UUID) (*vacancy.Vacancy, error)
	Publish(ctx context.Context, vacID uuid.UUID, compID uuid.UUID) error
}

type CompanyRepo interface {
	IncrementOpenVacancies(ctx context.Context, id uuid.UUID) error
}

type CacheRepo interface {
	Del(ctx context.Context, id uuid.UUID)
}

type CompMemberRepo interface {
	Get(ctx context.Context, userID, companyID uuid.UUID) (*domain.CompanyMember, error)
}
