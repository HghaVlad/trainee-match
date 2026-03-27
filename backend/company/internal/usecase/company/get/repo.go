package get

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
)

type CompanyRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*company.Company, error)
}

type CacheRepo interface {
	Get(ctx context.Context, key uuid.UUID) *company.Company
	Put(ctx context.Context, key uuid.UUID, val *company.Company, exp time.Duration)
}
