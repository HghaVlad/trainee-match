package get_company

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
)

type CompanyRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Company, error)
}

type CacheRepo interface {
	Get(ctx context.Context, key uuid.UUID) *domain.Company
	Put(ctx context.Context, key uuid.UUID, val *domain.Company, exp time.Duration)
}
