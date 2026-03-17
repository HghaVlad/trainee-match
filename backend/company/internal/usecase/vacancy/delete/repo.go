package delete_vacancy

import (
	"context"

	"github.com/google/uuid"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
)

type VacancyRepo interface {
	GetByID(ctx context.Context, vacancyID uuid.UUID, companyID uuid.UUID) (*domain.Vacancy, error)
	Delete(ctx context.Context, vacancyID uuid.UUID, companyID uuid.UUID) error
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
