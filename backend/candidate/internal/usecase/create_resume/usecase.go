package create_resume

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/google/uuid"
)

//go:generate mockery --name=ResumeRepo --output=mocks --outpkg=mocks
type ResumeRepo interface {
	Create(ctx context.Context, resume *domain.Resume) (uuid.UUID, error)
}

//go:generate mockery --name=SkillRepo --output=mocks --outpkg=mocks
type SkillRepo interface {
	AreSkillsExist(ctx context.Context, ids []uuid.UUID) (bool, error)
}

//go:generate mockery --name=CandidateRepo --output=mocks --outpkg=mocks
type CandidateRepo interface {
	GetByUserID(ctx context.Context, id uuid.UUID) (domain.Candidate, error)
}

type UseCase struct {
	resumeRepo    ResumeRepo
	skillRepo     SkillRepo
	candidateRepo CandidateRepo
}

func New(resumeRepo ResumeRepo, skillRepo SkillRepo, candidateRepo CandidateRepo) *UseCase {
	return &UseCase{
		resumeRepo:    resumeRepo,
		skillRepo:     skillRepo,
		candidateRepo: candidateRepo,
	}
}

func (uc *UseCase) Execute(ctx context.Context, req Request) (Response, error) {
	candidate, err := uc.candidateRepo.GetByUserID(ctx, req.UserId)
	if err != nil {
		return Response{}, err
	}

	// Convert request data to domain model
	domainData := convertRequestDataToDomainData(req.Data)

	resume := &domain.Resume{
		CandidateId: candidate.ID,
		Name:        req.Name,
		Status:      req.Status,
		Data:        domainData,
	}

	if err := resume.Validate(); err != nil {
		return Response{}, err
	}

	// Check if all skills exist
	if len(req.Data.SkillsList) > 0 {
		ok, err := uc.skillRepo.AreSkillsExist(ctx, req.Data.SkillsList)
		if err != nil {
			return Response{}, err
		}
		if !ok {
			return Response{}, domain.ErrSkillNotFound
		}
	}

	id, err := uc.resumeRepo.Create(ctx, resume)
	if err != nil {
		return Response{}, err
	}

	return Response{ID: id, CandidateID: candidate.ID}, nil
}

func convertRequestDataToDomainData(reqData ResumeData) domain.ResumeData {
	domainData := domain.ResumeData{
		LastName:        reqData.LastName,
		FirstName:       reqData.FirstName,
		MiddleName:      reqData.MiddleName,
		DateOfBirth:     reqData.DateOfBirth,
		Email:           reqData.Email,
		Phone:           reqData.Phone,
		City:            reqData.City,
		Citizenship:     reqData.Citizenship,
		Education:       make([]domain.Education, len(reqData.Education)),
		WorkExperiences: make([]domain.WorkExperience, len(reqData.WorkExperiences)),
		SkillsList:      reqData.SkillsList,
		AdditionalInfo:  reqData.AdditionalInfo,
		PortfolioLink:   reqData.PortfolioLink,
		DesiredFormat:   reqData.DesiredFormat,
		EnglishLevel:    reqData.EnglishLevel,
	}

	for i, edu := range reqData.Education {
		domainData.Education[i] = domain.Education{
			Level:          edu.Level,
			University:     edu.University,
			Faculty:        edu.Faculty,
			Specialization: edu.Specialization,
			StartYear:      edu.StartYear,
			EndYear:        edu.EndYear,
			Format:         edu.Format,
		}
	}

	for i, exp := range reqData.WorkExperiences {
		domainData.WorkExperiences[i] = domain.WorkExperience{
			Position:         exp.Position,
			Company:          exp.Company,
			Period:           exp.Period,
			Responsibilities: exp.Responsibilities,
		}
	}

	return domainData
}
