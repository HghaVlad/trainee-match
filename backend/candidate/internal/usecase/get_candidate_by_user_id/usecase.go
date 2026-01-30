package get_candidate_by_user_id

import (
	"context"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/google/uuid"
)

type CandidateRepo interface {
	GetByUserID(ctx context.Context, id uuid.UUID) (domain.Candidate, error)
}

type UseCase struct {
	repo CandidateRepo
}

func New(repo CandidateRepo) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Execute(ctx context.Context, id uuid.UUID) (*CandidateResponse, error) {
	candidate, err := uc.repo.GetByUserID(ctx, id)

	if err != nil {
		return nil, err
	}

	resp := CandidateResponse{
		ID:       candidate.ID,
		UserID:   candidate.UserId,
		Phone:    candidate.Phone,
		Telegram: candidate.Telegram,
		City:     candidate.City,
		Birthday: candidate.Birthday,
	}

	return &resp, nil
}
