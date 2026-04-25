package listbycomp

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type VacancyRepo interface {
	ListByCompanyByPublishedAt(ctx context.Context, compID uuid.UUID, cursor *PublishedAtCursor, limit int,
	) ([]VacancySummary, error)
}

type CompanyRepo interface {
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

type ResponseCacheRepo interface {
	Get(ctx context.Context, key string) *Response
	Put(ctx context.Context, key string, response *Response, exp time.Duration)
}
