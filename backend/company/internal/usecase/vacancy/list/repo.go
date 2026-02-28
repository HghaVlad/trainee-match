package list_vacancy

import (
	"context"
)

type VacancyRepo interface {
	List(ctx context.Context, requirements *Requirements, order Order, cursor any, limit int) ([]VacancySummary, error)
}
