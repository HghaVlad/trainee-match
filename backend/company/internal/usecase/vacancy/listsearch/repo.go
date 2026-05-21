package listsearch

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/views"
)

type VacancyRepo interface {
	ListPublishedSummaries(
		ctx context.Context,
		requirements *Requirements,
		order Order,
		cursor any,
		limit int,
	) (*SearchResult, error)
}

type SearchResult struct {
	Vacancies  []views.PublishedVacSummary
	NextCursor any
	HasNext    bool
}
