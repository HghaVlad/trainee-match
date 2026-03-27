package list

import (
	"context"
	"time"
)

type Repo interface {
	ListSummaries(ctx context.Context, order Order, cursor any, limit int) ([]CompanySummary, error)
}

type ResponseCacheRepo interface {
	Get(ctx context.Context, key string) *Response
	Put(ctx context.Context, key string, response *Response, exp time.Duration)
}
