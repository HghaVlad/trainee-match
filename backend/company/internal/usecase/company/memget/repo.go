package memget

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
)

type companyRepo interface {
	GetByMember(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*company.Company, error)
}
