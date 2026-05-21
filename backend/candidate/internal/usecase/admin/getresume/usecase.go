package getresume

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

type ResumeRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (domain.Resume, error)
}

type UseCase struct {
	repo ResumeRepo
}

func NewUseCase(repo ResumeRepo) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Execute(ctx context.Context, request Request) (Response, error) {
	resume, err := uc.repo.GetById(ctx, request.ResumeID)

	if err != nil {
		return Response{}, err
	}

	status, err := domain.Format(resume.Status)
	if err != nil {
		return Response{}, err
	}

	response := Response{
		ID:               resume.ID,
		CandidateID:      resume.CandidateId,
		Name:             resume.Name,
		Status:           string(status),
		ModerationStatus: string(resume.ModerationStatus),
		Data:             convertDomainDataToResponseData(resume.Data),
	}

	return response, nil
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
