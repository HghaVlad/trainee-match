package list

import (
	"context"
)

type VacancyRepo interface {
	ListPublished(
		ctx context.Context,
		requirements *Requirements,
		order Order,
		cursor any,
		limit int,
	) ([]VacancySummary, error)
}
