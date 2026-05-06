package update

import (
	"context"

	"github.com/google/uuid"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
)

//go:generate mockgen -source=repo.go -destination=mocks/repo_mocks.go -package=mocks
type CompanyRepo interface {
	UpdateAndGetOldName(ctx context.Context, req *Request) (string, error)
}

type CompMemberRepo interface {
	Get(ctx context.Context, userID, companyID uuid.UUID) (*domain.CompanyMember, error)
}

type CacheRepo interface {
	Del(ctx context.Context, id uuid.UUID)
}
