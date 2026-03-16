package get_vacancy

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID, companyID uuid.UUID) (*domain.Vacancy, error)
}

type CacheRepo interface {
	Get(ctx context.Context, key uuid.UUID) *domain.Vacancy
	Put(ctx context.Context, key uuid.UUID, val *domain.Vacancy, exp time.Duration)
}

type CompMemberRepo interface {
	Get(ctx context.Context, userID, companyID uuid.UUID) (*domain.CompanyMember, error)
}
