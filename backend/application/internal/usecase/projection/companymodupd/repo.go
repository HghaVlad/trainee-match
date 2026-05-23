package companymodupd

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type vacancyRepo interface {
	UpdateCompanyModStatus(ctx context.Context, compID uuid.UUID, status projection.ModerationStatus) error
}
