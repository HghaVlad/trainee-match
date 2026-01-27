package list_companies

import (
	"context"
)

type Repo interface {
	ListByCreatedAtDesc(ctx context.Context, cursor *CreatedAtCursor, limit int) (
		[]CompanySummary, *CreatedAtCursor, error)

	ListByName(ctx context.Context, cursor *NameCursor, limit int) (
		[]CompanySummary, *NameCursor, error)

	ListByVacanciesCnt(ctx context.Context, cursor *VacanciesCntCursor, limit int) (
		[]CompanySummary, *VacanciesCntCursor, error)
}
