package get

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID, companyID uuid.UUID) (*vacancy.Vacancy, error)
}

type CacheRepo interface {
	Get(ctx context.Context, key uuid.UUID) *vacancy.Vacancy
	Put(ctx context.Context, key uuid.UUID, val *vacancy.Vacancy, exp time.Duration)
}

type CompMemberRepo interface {
	Get(ctx context.Context, userID, companyID uuid.UUID) (*member.CompanyMember, error)
}
