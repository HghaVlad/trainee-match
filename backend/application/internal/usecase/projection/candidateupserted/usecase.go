package candidateupserted

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type Usecase struct {
	repo CandidateRepo
}

func NewUsecase(repo CandidateRepo) *Usecase {
	return &Usecase{repo: repo}
}

func (uc *Usecase) Execute(ctx context.Context, event projection.CandidateUpsertedEvent) error {
	candidate := event.ToCandidate()
	return uc.repo.Save(ctx, &candidate)
}
