package get_company

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
)

type CompanyRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Company, error)
}
