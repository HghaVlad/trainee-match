package add

import (
	"context"

	"github.com/google/uuid"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/projection/userhr"
)

//go:generate mockgen -source=repo.go -destination=mocks/repo_mocks.go -package=mocks
type companyMemberRepo interface {
	Get(ctx context.Context, userID, companyID uuid.UUID) (*domain.CompanyMember, error)
	Create(ctx context.Context, member *domain.CompanyMember) error
}

type hrProjRepo interface {
	GetByUsername(ctx context.Context, username string) (*userhr.Projection, error)
}
