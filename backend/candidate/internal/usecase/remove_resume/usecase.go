package remove_resume

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain/events"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

//go:generate mockery --name=ResumeRepo --output=mocks --outpkg=mocks
type ResumeRepo interface {
	Remove(ctx context.Context, resumeId uuid.UUID, candidateId uuid.UUID) error
}

//go:generate mockery --name=CandidateRepo --output=mocks --outpkg=mocks
type CandidateRepo interface {
	GetByUserID(ctx context.Context, id uuid.UUID) (domain.Candidate, error)
}

type EventWriter interface {
	WriteResumeDeleted(ctx context.Context, ev events.ResumeDeleted) error
}

type Usecase struct {
	resumeRepo    ResumeRepo
	candidateRepo CandidateRepo
	writer        EventWriter
	trManager     *manager.Manager
}

func New(resumeRepo ResumeRepo, candidateRepo CandidateRepo, writer EventWriter, trManager *manager.Manager) *Usecase {
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
