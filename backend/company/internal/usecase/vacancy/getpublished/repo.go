package getpublished

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	GetPublishedByID(ctx context.Context, id uuid.UUID) (*Response, error)
}

type CacheRepo interface {
	Get(ctx context.Context, key uuid.UUID) *Response
	Put(ctx context.Context, key uuid.UUID, val *Response, exp time.Duration)
}
