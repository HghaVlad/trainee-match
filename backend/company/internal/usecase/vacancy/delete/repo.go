package delete_vacancy

import (
	"context"

	"github.com/google/uuid"
)

type VacancyRepo interface {
	Delete(ctx context.Context, vacancyID uuid.UUID, companyID uuid.UUID) error
}

type CacheRepo interface {
	Del(ctx context.Context, id uuid.UUID)
}
