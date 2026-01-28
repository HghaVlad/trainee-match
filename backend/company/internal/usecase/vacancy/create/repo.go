package create_vacancy

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
)

type VacancyRepo interface {
	Create(ctx context.Context, vacancy *domain.Vacancy) error
}

type CompanyRepo interface {
	IncrementOpenVacancies(ctx context.Context, id uuid.UUID) error
}
