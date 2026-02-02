package list_vacancy

import (
	"context"
	"time"
)

type VacancyRepo interface {
	ListByPublishedAt(ctx context.Context, cursor *PublishedAtCursor, limit int) (
		[]VacancySummary, *PublishedAtCursor, error)
}

type ResponseCacheRepo interface {
	Get(ctx context.Context, key string) *Response
	Put(ctx context.Context, key string, response *Response, exp time.Duration)
}
