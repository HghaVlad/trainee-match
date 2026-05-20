package vacancypublished

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type VacancyRepo interface {
	Save(ctx context.Context, vacancy *projection.Vacancy) error
}
