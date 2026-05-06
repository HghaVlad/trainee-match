package list

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
)

type repo interface {
	Get(ctx context.Context, userID, companyID uuid.UUID) (*member.CompanyMember, error)
	ListViewsByCompany(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]View, error)
}
