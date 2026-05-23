package getcandidate

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

type CandidateRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (domain.Candidate, error)
}

type UseCase struct {
	repo CandidateRepo
}

func NewUseCase(repo CandidateRepo) *UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) Execute(ctx context.Context, request Request) (CandidateResponse, error) {
	candidate, err := u.repo.GetByID(ctx, request.CandidateID)
	if err != nil {
		return CandidateResponse{}, err
	}

	return CandidateResponse{
		ID:       candidate.ID,
		FullName: candidate.FullName,
		UserID:   candidate.UserId,
		Phone:    candidate.Phone,
		Telegram: candidate.Telegram,
		City:     candidate.City,
		Birthday: candidate.Birthday,
	}, nil
}
