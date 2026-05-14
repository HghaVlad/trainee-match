package resumedeleted

import (
	"context"

	"github.com/google/uuid"
)

type ResumeRepo interface {
	Delete(ctx context.Context, resumeID uuid.UUID) error
}
