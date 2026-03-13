package publish_vacancy

import (
	"context"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/google/uuid"
)

type VacancyRepo interface {
	Publish(ctx context.Context, compID uuid.UUID, vacID uuid.UUID) error
}

type CompMemberRepo interface {
	Get(ctx context.Context, userID, companyID uuid.UUID) (*domain.CompanyMember, error)
}
