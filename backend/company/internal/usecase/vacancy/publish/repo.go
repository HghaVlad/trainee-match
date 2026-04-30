package publish

import (
	"context"

	"github.com/google/uuid"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
)

//go:generate mockgen -source=repo.go -destination=mocks/repo_mocks.go -package=mocks
type VacancyRepo interface {
	GetPublishedEventView(
		ctx context.Context,
		vacancyID uuid.UUID,
		companyID uuid.UUID,
	) (*PublishedEventView, error)
	Publish(ctx context.Context, vacID uuid.UUID, compID uuid.UUID) error
}

type CompanyRepo interface {
	IncrementOpenVacancies(ctx context.Context, id uuid.UUID) error
}

type CacheRepo interface {
	Del(ctx context.Context, id uuid.UUID)
}

type CompMemberRepo interface {
	Get(ctx context.Context, userID, companyID uuid.UUID) (*domain.CompanyMember, error)
}
