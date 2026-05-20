package list

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
	) ([]views.PublishedVacSummary, error)
}
