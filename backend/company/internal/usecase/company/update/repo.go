package update_company

import (
	"context"

	"github.com/google/uuid"
)

type CompanyRepo interface {
	Update(ctx context.Context, req *Request) error
}

type CacheRepo interface {
	Del(ctx context.Context, id uuid.UUID)
}
