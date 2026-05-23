package deleteskill

import (
	"context"

	"github.com/google/uuid"
)

type SkillRepo interface {
	Delete(ctx context.Context, id uuid.UUID) error
}

type UseCase struct {
	repo SkillRepo
}

func NewUseCase(repo SkillRepo) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Execute(ctx context.Context, req Request) error {
	return uc.repo.Delete(ctx, req.SkillID)
}
