package create_candidate

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain/events"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

//go:generate mockery
type CandidateRepo interface {
	Create(ctx context.Context, candidate *domain.Candidate) (uuid.UUID, error)
	GetByUserID(ctx context.Context, id uuid.UUID) (domain.Candidate, error)
}

type EventWriter interface {
	WriteCandidateUpserted(ctx context.Context, ev events.CandidateUpserted) error
}

type TrManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type UseCase struct {
	repo      CandidateRepo
	writer    EventWriter
	trManager TrManager
}

func New(repo CandidateRepo, writer EventWriter, trManager TrManager) *UseCase {
	return &UseCase{repo: repo, writer: writer, trManager: trManager}
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

	err := uc.trManager.Do(ctx, func(ctx context.Context) error {
		id, err := uc.repo.Create(ctx, candidate)
		if err != nil {
			return err
		}
		candidate.ID = id
		return uc.writer.WriteCandidateUpserted(ctx, events.NewCandidateUpserted(*candidate, req.FullName, req.Email))
	})
	if err != nil {
		return uuid.Nil, err
	}

	return candidate.ID, nil

}
