package create_candidate

import (
	"context"
	"errors"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/google/uuid"
)

//go:generate mockery --name=CandidateRepo --output=mocks --outpkg=mocks
type CandidateRepo interface {
	Create(ctx context.Context, candidate *domain.Candidate) (uuid.UUID, error)
	GetByUserID(ctx context.Context, id uuid.UUID) (domain.Candidate, error)
}

type UseCase struct {
	repo CandidateRepo
}

func New(repo CandidateRepo) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Execute(ctx context.Context, req *Request) (uuid.UUID, error) {
	if _, err := uc.repo.GetByUserID(ctx, req.UserID); err == nil {
		return uuid.Nil, domain.ErrCandidateAlreadyExists
	} else if !errors.Is(err, domain.ErrCandidateNotFound) {
		return uuid.Nil, err
	}

	candidate := &domain.Candidate{
		UserId:   req.UserID,
		Phone:    req.Phone,
		Telegram: req.Telegram,
		City:     req.City,
		Birthday: req.Birthday,
	}
	if err := candidate.Validate(); err != nil {
		return uuid.Nil, err
	}

	id, err := uc.repo.Create(ctx, candidate)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil

}
