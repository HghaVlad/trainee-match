package vacancyarchived

import (
	"context"

	"github.com/google/uuid"
)

type VacancyRepo interface {
	Archive(ctx context.Context, vacancyID uuid.UUID) error
}
