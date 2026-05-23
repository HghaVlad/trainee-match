package addskill

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

type SkillRepo interface {
	Create(ctx context.Context, skill domain.Skill) (uuid.UUID, error)
}

type UseCase struct {
	repo SkillRepo
}

func NewUseCase(repo SkillRepo) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Execute(ctx context.Context, req Request) (Response, error) {
	skill := domain.Skill{Name: req.Name}
	if err := skill.Validate(); err != nil {
		return Response{}, err
	}

	id, err := uc.repo.Create(ctx, skill)
	if err != nil {
		return Response{}, err
	}

	return Response{ID: id, Name: skill.Name}, nil
}
