package vacancymodupd

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type vacancyRepo interface {
	UpdateModStatus(ctx context.Context, vacID uuid.UUID, status projection.ModerationStatus) error
}
