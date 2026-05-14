package candidateupserted

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type CandidateRepo interface {
	Save(ctx context.Context, candidate *projection.Candidate) error
}
