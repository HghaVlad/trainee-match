package getcandidates

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

type CandidateRepo interface {
	GetCandidates(ctx context.Context, page int, size int) ([]domain.Candidate, error)
}

type UseCase struct {
	repo CandidateRepo
}

func NewUseCase(repo CandidateRepo) *UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) Execute(ctx context.Context, req Request) ([]CandidateResponse, error) {
	candidates, err := u.repo.GetCandidates(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	resp := make([]CandidateResponse, len(candidates))
	for i, c := range candidates {
		resp[i] = CandidateResponse{
			ID:       c.ID,
			FullName: c.FullName,
			UserID:   c.UserId,
			Phone:    c.Phone,
			Telegram: c.Telegram,
			City:     c.City,
			Birthday: c.Birthday,
		}
	}
	return resp, nil
}
