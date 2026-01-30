package create_candidate

import (
	"context"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/google/uuid"
)

type CandidateRepo interface {
	Create(ctx context.Context, candidate *domain.Candidate) (uuid.UUID, error)
}

type UseCase struct {
	repo CandidateRepo
}

func New(repo CandidateRepo) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Execute(ctx context.Context, req *Request) (uuid.UUID, error) {
	candidate := &domain.Candidate{
		UserId:   req.UserID,
		Phone:    req.Phone,
		Telegram: req.Telegram,
		City:     req.City,
		Birthday: req.Birthday,
	}

	if id, err := uc.repo.Create(ctx, candidate); err != nil {
		return uuid.Nil, err
	} else {
		return id, nil
	}
}
