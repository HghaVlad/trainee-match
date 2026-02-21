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
	GetByTelegram(ctx context.Context, telegram string) (domain.Candidate, error)
	GetByPhone(ctx context.Context, phone string) (domain.Candidate, error)
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
	if err := candidate.Validate(); err != nil {
		return uuid.Nil, err
	}
	if err := uc.validateUniqueness(ctx, candidate); err != nil {
		return uuid.Nil, err
	}

	if id, err := uc.repo.Create(ctx, candidate); err != nil {
		return uuid.Nil, err
	} else {
		return id, nil
	}
}

func (uc *UseCase) validateUniqueness(ctx context.Context, candidate *domain.Candidate) error {
	_, err := uc.repo.GetByTelegram(ctx, candidate.Telegram)
	if err == nil {
		return domain.ErrTelegramAlreadyExists
	} else if !errors.Is(err, domain.ErrCandidateNotFound) {
		return err
	}

	_, err = uc.repo.GetByPhone(ctx, candidate.Phone)
	if err == nil {
		return domain.ErrPhoneAlreadyExists
	} else if !errors.Is(err, domain.ErrCandidateNotFound) {
		return err
	}

	return nil
}
