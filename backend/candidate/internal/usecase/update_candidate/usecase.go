package update_candidate

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/google/uuid"
)

//go:generate mockery --name=CandidateRepo --output=mocks --outpkg=mocks
type CandidateRepo interface {
	Update(ctx context.Context, candidate domain.Candidate) (domain.Candidate, error)
	GetByUserID(ctx context.Context, id uuid.UUID) (domain.Candidate, error)
}

type UseCase struct {
	repo CandidateRepo
}

func New(repo CandidateRepo) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Execute(ctx context.Context, userID uuid.UUID, req *Request) (*CandidateResponse, error) {
	candidate, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if candidate.UserId != userID {
		return nil, domain.ErrForbidden
	}

	if req.UserID != nil {
		// prevent changing owner to another user
		if *req.UserID != userID {
			return nil, domain.ErrForbidden
		}
		candidate.UserId = *req.UserID
	}
	if req.Phone != nil {
		candidate.Phone = *req.Phone
	}
	if req.Telegram != nil {
		candidate.Telegram = *req.Telegram
	}
	if req.City != nil {
		candidate.City = *req.City
	}
	if req.Birthday != nil {
		candidate.Birthday = *req.Birthday
	}

	if err = candidate.Validate(); err != nil {
		return nil, err
	}

	candidate, err = uc.repo.Update(ctx, candidate)
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
