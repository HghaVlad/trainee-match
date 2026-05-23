package remove_resume

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain/events"
)

type ResumeRepo interface {
	Remove(ctx context.Context, resumeId uuid.UUID, candidateId uuid.UUID) error
}

type CandidateRepo interface {
	GetByUserID(ctx context.Context, id uuid.UUID) (domain.Candidate, error)
}

type EventWriter interface {
	WriteResumeDeleted(ctx context.Context, ev events.ResumeDeleted) error
}

type TrManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type Usecase struct {
	resumeRepo    ResumeRepo
	candidateRepo CandidateRepo
	writer        EventWriter
	trManager     TrManager
}

func New(resumeRepo ResumeRepo, candidateRepo CandidateRepo, writer EventWriter, trManager TrManager) *Usecase {
	return &Usecase{resumeRepo: resumeRepo, candidateRepo: candidateRepo, writer: writer, trManager: trManager}
}

func (uc *Usecase) RemoveResume(ctx context.Context, req Request) error {
	candidate, err := uc.candidateRepo.GetByUserID(ctx, req.UserId)
	if err != nil {
		return err
	}

	err = uc.trManager.Do(ctx, func(ctx context.Context) error {
		err = uc.resumeRepo.Remove(ctx, req.ResumeId, candidate.ID)
		if err != nil {
			return err
		}
		return uc.writer.WriteResumeDeleted(ctx, events.NewResumeDeleted(req.ResumeId))
	})

	return err
}
