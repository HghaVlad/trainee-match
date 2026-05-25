package update

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

//go:generate mockgen -source=repo.go -destination=mocks/repo_mocks.go -package=mocks
type VacancyRepo interface {
	GetByIDForUpdate(ctx context.Context, vacancyID uuid.UUID, companyID uuid.UUID) (*vacancy.Vacancy, error)
	Update(ctx context.Context, v *vacancy.Vacancy) error
}

type compRepo interface {
	GetByMember(ctx context.Context, compID, userID uuid.UUID) (*company.Company, error)
}

type CacheRepo interface {
	Del(ctx context.Context, id uuid.UUID)
}
