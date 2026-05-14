package resumeupserted

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type ResumeRepo interface {
	Save(ctx context.Context, resume *projection.Resume) error
}
