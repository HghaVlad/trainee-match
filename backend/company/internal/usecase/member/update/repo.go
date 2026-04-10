package update

import (
	"context"

	"github.com/google/uuid"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
)

type CompanyMemberRepo interface {
	Get(ctx context.Context, userID, companyID uuid.UUID) (*domain.CompanyMember, error)
	UpdateRole(ctx context.Context, userID, companyID uuid.UUID, role domain.CompanyRole) error
}
