package get_vacancy

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID, companyID uuid.UUID) (*domain.Vacancy, error)
}
