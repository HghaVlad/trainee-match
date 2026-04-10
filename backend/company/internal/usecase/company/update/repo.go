package update

import (
	"context"

	"github.com/google/uuid"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
)

type CompanyRepo interface {
	Update(ctx context.Context, req *Request) error
}

type CompMemberRepo interface {
	Get(ctx context.Context, userID, companyID uuid.UUID) (*domain.CompanyMember, error)
}

type CacheRepo interface {
	Del(ctx context.Context, id uuid.UUID)
}
