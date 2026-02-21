package update_candidate

import (
	"context"
	"errors"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/google/uuid"
)

//go:generate mockery --name=CandidateRepo --output=mocks --outpkg=mocks
type CandidateRepo interface {
	Update(ctx context.Context, candidate domain.Candidate) (domain.Candidate, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Candidate, error)
	GetByTelegram(ctx context.Context, telegram string) (domain.Candidate, error)
	GetByPhone(ctx context.Context, phone string) (domain.Candidate, error)
}

type UseCase struct {
	repo CandidateRepo
}

func New(repo CandidateRepo) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Execute(ctx context.Context, req *Request) (*CandidateResponse, error) {
	candidate, err := uc.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if req.UserID != nil {
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

	if err = uc.validateUniqueness(ctx, candidate); err != nil {
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

func (uc *UseCase) validateUniqueness(ctx context.Context, candidate domain.Candidate) error {
	other, err := uc.repo.GetByTelegram(ctx, candidate.Telegram)
	if err == nil && other.ID != candidate.ID {
		return domain.ErrTelegramAlreadyExists
	}
	if err != nil && !errors.Is(domain.ErrCandidateNotFound, err) {
		return err
	}

	other, err = uc.repo.GetByPhone(ctx, candidate.Phone)
	if err == nil && other.ID != candidate.ID {
		return domain.ErrPhoneAlreadyExists
	}
	if err != nil && !errors.Is(domain.ErrCandidateNotFound, err) {
		return err
	}

	return nil
}
