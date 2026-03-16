package archive_vacancy

import (
	"context"

	"github.com/google/uuid"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
)

type VacancyRepo interface {
	GetByID(ctx context.Context, vacID uuid.UUID, compID uuid.UUID) (*domain.Vacancy, error)
	Archive(ctx context.Context, compID uuid.UUID, vacID uuid.UUID) error
}

type CompanyRepo interface {
	DecrementOpenVacancies(ctx context.Context, id uuid.UUID) error
}

type CompMemberRepo interface {
	Get(ctx context.Context, userID, companyID uuid.UUID) (*domain.CompanyMember, error)
}

type CacheRepo interface {
	Del(ctx context.Context, id uuid.UUID)
}
