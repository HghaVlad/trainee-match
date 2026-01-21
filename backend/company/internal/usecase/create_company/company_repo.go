package create_company

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
)

type CompanyRepo interface {
	Create(ctx context.Context, company *entities.Company) error
}
