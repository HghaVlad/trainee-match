package vacancyupdated

import (
	"context"

	"github.com/google/uuid"
)

type VacancyRepo interface {
	UpdateTitle(ctx context.Context, vacancyID uuid.UUID, title string) error
}
