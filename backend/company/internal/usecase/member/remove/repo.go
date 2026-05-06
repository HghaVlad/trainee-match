package remove

import (
	"context"

	"github.com/google/uuid"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
)

//go:generate mockgen -source=repo.go -destination=mocks/repo_mocks.go -package=mocks
type CompanyMemberRepo interface {
	Get(ctx context.Context, userID, companyID uuid.UUID) (*domain.CompanyMember, error)
	GetCompanyRoleCount(ctx context.Context, companyID uuid.UUID, role domain.CompanyRole) (int, error)
	Delete(ctx context.Context, userID, companyID uuid.UUID) error
}
