package getcandidateresumes

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

type ResumeRepo interface {
	GetByCandidateId(
		ctx context.Context,
		candidateID uuid.UUID,
		page, size int,
	) (resumes []domain.Resume, err error)
}

type UseCase struct {
	repo ResumeRepo
}

func NewUseCase(repo ResumeRepo) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Execute(ctx context.Context, request Request) ([]ShortResponse, error) {
	resumes, err := uc.repo.GetByCandidateId(ctx, request.CandidateID, request.Page, request.Size)
	if err != nil {
		return nil, err
	}

	resp := make([]ShortResponse, len(resumes))
	for i, resume := range resumes {
		status, err := domain.Format(resume.Status)
		if err != nil {
			return nil, err
		}
		resp[i] = ShortResponse{
			ID:               resume.ID,
			Name:             resume.Name,
			CandidateID:      resume.CandidateId,
			Status:           string(status),
			ModerationStatus: string(resume.ModerationStatus),
		}
	}
	return resp, nil
}
