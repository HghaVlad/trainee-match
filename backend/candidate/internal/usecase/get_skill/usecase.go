package get_skill

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/google/uuid"
)

//go:generate mockery --name=SkillRepo --output=mocks --outpkg=mocks
type SkillRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (domain.Skill, error)
	List(ctx context.Context) ([]domain.Skill, error)
}

type UseCase struct {
	repo SkillRepo
}

func New(repo SkillRepo) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Execute(ctx context.Context, req GetByIdRequest) (*GetByIdResponse, error) {
	skill, err := uc.repo.GetByID(ctx, req.ID)

	if err != nil {
		return nil, err
	}

	response := &GetByIdResponse{
		ID:   skill.ID,
		Name: skill.Name,
	}

	return response, nil
}

// TODO: implement pagination
func (uc *UseCase) ExecuteList(ctx context.Context, req ListRequest) ([]*ListResponse, error) {
	skills, err := uc.repo.List(ctx)

	if err != nil {
		return nil, err
	}

	var result []*ListResponse
	for _, skill := range skills {
		if err := skill.Validate(); err != nil {
			return nil, err
		}
		item := &ListResponse{
			ID:   skill.ID,
			Name: skill.Name,
		}
		result = append(result, item)
	}

	return result, nil
}
