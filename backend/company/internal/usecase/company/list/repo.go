package list

import (
	"context"
	"time"
)

type Repo interface {
	ListByCreatedAtDesc(ctx context.Context, cursor *CreatedAtCursor, limit int) (
		[]CompanySummary, *CreatedAtCursor, error)

	ListByName(ctx context.Context, cursor *NameCursor, limit int) (
		[]CompanySummary, *NameCursor, error)

	ListByVacanciesCnt(ctx context.Context, cursor *VacanciesCntCursor, limit int) (
		[]CompanySummary, *VacanciesCntCursor, error)
}

type ResponseCacheRepo interface {
	Get(ctx context.Context, key string) *Response
	Put(ctx context.Context, key string, response *Response, exp time.Duration)
}
