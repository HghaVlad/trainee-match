package publish

import (
	"context"

	"github.com/google/uuid"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
)

//go:generate mockgen -source=repo.go -destination=mocks/repo_mocks.go -package=mocks
type VacancyRepo interface {
	// PublishIfNotPublished if vacancy exists and is not published,
	// marks it as published and returns view for the event.
	// If it was already published, updates nothing and returns view WasAlreadyPublished = true.
	// If vacancy doesn't exist, returns vacancy.ErrVacancyNotFound.
	// Optimized to do a single round trip to db to be effective and avoid race conditions.
	PublishIfNotPublished(ctx context.Context, vacID, compID uuid.UUID) (*PublishedEventView, error)
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
