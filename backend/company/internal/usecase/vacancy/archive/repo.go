package archive_vacancy

import (
	"context"
	"time"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
	"github.com/google/uuid"
)

type VacancyRepo interface {
	UpdateStatus(ctx context.Context, compID uuid.UUID, vacID uuid.UUID,
		status value_types.VacancyStatus, pubTime *time.Time) error
}

type CompMemberRepo interface {
	Get(ctx context.Context, userID, companyID uuid.UUID) (*domain.CompanyMember, error)
}
