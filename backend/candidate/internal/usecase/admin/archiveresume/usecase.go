package archiveresume

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain/events"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

type ResumeRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (domain.Resume, error)
	SetModerationStatus(ctx context.Context, id uuid.UUID, status domain.ModerationStatus) error
}

type EventWriter interface {
	WriteResumeArchived(ctx context.Context, ev events.ResumeArchived) error
}

type TrManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type UseCase struct {
	repo      ResumeRepo
	writer    EventWriter
	trManager TrManager
}

func NewUseCase(repo ResumeRepo, writer EventWriter, manager TrManager) *UseCase {
	return &UseCase{repo: repo, writer: writer, trManager: manager}
}

func (u *UseCase) Execute(ctx context.Context, req Request) error {
	resume, err := u.repo.GetById(ctx, req.ResumeID)
	if err != nil {
		return err
	}

	if resume.ModerationStatus == domain.ModerationStatusHidden {
		return nil
	}

	err = u.trManager.Do(ctx, func(ctx context.Context) error {
		err = u.repo.SetModerationStatus(ctx, req.ResumeID, domain.ModerationStatusHidden)
		if err != nil {
			return err
		}
		return u.writer.WriteResumeArchived(ctx, events.NewResumeArchived(resume.ID))
	})

	return err
}
