package resumearchived

import (
	"context"

	"github.com/google/uuid"
)

type ResumeRepo interface {
	Archive(ctx context.Context, resumeID uuid.UUID) error
}
