package update_candidate

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain/events"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

type CandidateRepo interface {
	Update(ctx context.Context, candidate domain.Candidate) (domain.Candidate, error)
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
	if req.FullName != "" {
		candidate.FullName = req.FullName
	}

	if err = candidate.Validate(); err != nil {
		return nil, err
	}

	err = uc.trManager.Do(ctx, func(ctx context.Context) error {
		candidate, err = uc.repo.Update(ctx, candidate)
		if err != nil {
			return err
		}
		return uc.writer.WriteCandidateUpserted(ctx, events.NewCandidateUpserted(candidate, req.Email))
	})
	if err != nil {
		return nil, err
	}

	resp := CandidateResponse{
		ID:       candidate.ID,
		UserID:   candidate.UserId,
		FullName: candidate.FullName,
		Phone:    candidate.Phone,
		Telegram: candidate.Telegram,
		City:     candidate.City,
		Birthday: candidate.Birthday,
	}

	return &resp, nil
}
