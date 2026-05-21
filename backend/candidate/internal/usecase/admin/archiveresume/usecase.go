package archiveresume

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

type ResumeRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (domain.Resume, error)
	SetModerationStatus(ctx context.Context, id uuid.UUID, status domain.ModerationStatus) error
}

type UseCase struct {
	repo ResumeRepo
}

func NewUseCase(repo ResumeRepo) *UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) Execute(ctx context.Context, req Request) error {
	resume, err := u.repo.GetById(ctx, req.ResumeID)
	if err != nil {
		return err
	}

	if resume.ModerationStatus == domain.ModerationStatusHidden {
		return nil
	}

	return u.repo.SetModerationStatus(ctx, req.ResumeID, domain.ModerationStatusHidden)
}
