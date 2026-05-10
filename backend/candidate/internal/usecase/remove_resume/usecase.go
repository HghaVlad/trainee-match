package remove_resume

import (
	"context"

	"github.com/google/uuid"

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

type Usecase struct {
	resumeRepo    ResumeRepo
	candidateRepo CandidateRepo
}

func New(resumeRepo ResumeRepo, candidateRepo CandidateRepo) *Usecase {
	return &Usecase{resumeRepo: resumeRepo, candidateRepo: candidateRepo}
}

func (uc *Usecase) RemoveResume(ctx context.Context, req Request) error {
	candidate, err := uc.candidateRepo.GetByUserID(ctx, req.UserId)
	if err != nil {
		return err
	}

	err = uc.resumeRepo.Remove(ctx, req.ResumeId, candidate.ID)
	return err
}
