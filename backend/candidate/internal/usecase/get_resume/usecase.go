package get_resume

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/google/uuid"
)

//go:generate mockery --name=ResumeRepo --output=mocks --outpkg=mocks
type ResumeRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (domain.Resume, error)
	GetByCandidateId(ctx context.Context, candidateId uuid.UUID) ([]domain.Resume, error)
}

//go:generate mockery --name=CandidateRepo --output=mocks --outpkg=mocks
type CandidateRepo interface {
	GetByUserID(ctx context.Context, id uuid.UUID) (domain.Candidate, error)
}

type UseCase struct {
	resumeRepo    ResumeRepo
	candidateRepo CandidateRepo
}

func New(resumeRepo ResumeRepo, candidateRepo CandidateRepo) *UseCase {
	return &UseCase{
		resumeRepo:    resumeRepo,
		candidateRepo: candidateRepo,
	}
}

// TODO: add check that user has access to this resume
func (uc *UseCase) GetById(ctx context.Context, resumeId, UserId uuid.UUID) (*Response, error) {
	resume, err := uc.resumeRepo.GetById(ctx, resumeId)

	if err != nil {
		return nil, err
	}

	response := &Response{
		ID:          resume.ID,
		CandidateID: resume.CandidateId,
		Name:        resume.Name,
		Status:      resume.Status,
		Data:        convertDomainDataToResponseData(resume.Data),
	}

	return response, nil
}

func (uc *UseCase) GetByCandidateId(ctx context.Context, UserId uuid.UUID) ([]*ShortResponse, error) {
	candidate, err := uc.candidateRepo.GetByUserID(ctx, UserId)
	if err != nil {
		return nil, err
	}
	resumes, err := uc.resumeRepo.GetByCandidateId(ctx, candidate.ID)

	if err != nil {
		return nil, err
	}

	var result []*ShortResponse
	for _, resume := range resumes {
		item := &ShortResponse{
			ID:          resume.ID,
			CandidateId: resume.CandidateId,
			Name:        resume.Name,
			Status:      resume.Status,
		}
		result = append(result, item)
	}

	return result, nil
}

// Helper function to convert domain data to response data
func convertDomainDataToResponseData(domainData domain.ResumeData) ResumeData {

	responseData := ResumeData{
		LastName:        domainData.LastName,
		FirstName:       domainData.FirstName,
		MiddleName:      domainData.MiddleName,
		DateOfBirth:     domainData.DateOfBirth,
		Email:           domainData.Email,
		Phone:           domainData.Phone,
		City:            domainData.City,
		Citizenship:     domainData.Citizenship,
		Education:       make([]Education, len(domainData.Education)),
		WorkExperiences: make([]WorkExperience, len(domainData.WorkExperiences)),
		SkillsList:      domainData.SkillsList,
		AdditionalInfo:  domainData.AdditionalInfo,
		PortfolioLink:   domainData.PortfolioLink,
		DesiredFormat:   domainData.DesiredFormat,
		EnglishLevel:    domainData.EnglishLevel,
	}

	for i, edu := range domainData.Education {
		responseData.Education[i] = Education{
			Level:          edu.Level,
			University:     edu.University,
			Faculty:        edu.Faculty,
			Specialization: edu.Specialization,
			StartYear:      edu.StartYear,
			EndYear:        edu.EndYear,
			Format:         edu.Format,
		}
	}

	for i, exp := range domainData.WorkExperiences {
		responseData.WorkExperiences[i] = WorkExperience{
			Position:         exp.Position,
			Company:          exp.Company,
			Period:           exp.Period,
			Responsibilities: exp.Responsibilities,
		}
	}

	return responseData
}
