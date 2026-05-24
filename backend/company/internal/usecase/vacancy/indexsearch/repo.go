package indexsearch

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/views"
)

type VacancyRepo interface {
	GetSearchView(ctx context.Context, vacID uuid.UUID) (*views.VacancySearch, error)
}

type SearchRepo interface {
	Index(ctx context.Context, vac views.VacancySearch) error
	UpdateCompanyName(ctx context.Context, companyID uuid.UUID, companyName string) error
	RemoveByCompanyID(ctx context.Context, compID uuid.UUID) error
}
