package archive

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/views"
)

//go:generate mockgen -source=repo.go -destination=mocks/repo_mocks.go -package=mocks
type VacancyRepo interface {
	ArchiveAndGetOldStatus(ctx context.Context, vacID, compID uuid.UUID) (vacancy.Status, error)
	GetSearchView(ctx context.Context, vacID uuid.UUID) (*views.VacancySearch, error)
}

type CompanyRepo interface {
	DecrementOpenVacancies(ctx context.Context, id uuid.UUID) error
}

type CompMemberRepo interface {
	Get(ctx context.Context, userID, companyID uuid.UUID) (*member.CompanyMember, error)
}

type CacheRepo interface {
	Del(ctx context.Context, id uuid.UUID)
}

type SearchRepo interface {
	Index(ctx context.Context, vac views.VacancySearch) error
}
