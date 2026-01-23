package delete_company

import (
	"context"

	"github.com/google/uuid"
)

type CompanyRepo interface {
	Delete(ctx context.Context, id uuid.UUID) error
}
