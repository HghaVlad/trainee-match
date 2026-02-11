package get_resume

import (
	"context"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/google/uuid"
)

type ResumeRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (domain.Resume, error)
	GetByCandidateId(ctx context.Context, candidateId uuid.UUID) ([]domain.Resume, error)
}

type UseCase struct {
	repo ResumeRepo
}

func New(repo ResumeRepo) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) GetById(ctx context.Context, req GetByIdRequest) (*GetByIdResponse, error) {
	resume, err := uc.repo.GetById(ctx, req.ID)

	if err != nil {
		return nil, err
	}

	// Convert domain model to response
	response := &GetByIdResponse{
		ID:          resume.ID,
		CandidateId: resume.CandidateId,
		Name:        resume.Name,
		Status:      resume.Status,
		Data:        convertDomainDataToResponseData(resume.Data),
	}

	return response, nil
}

func (uc *UseCase) GetByCandidateId(ctx context.Context, req GetByCandidateIdRequest) ([]*GetByCandidateIdResponse, error) {
	resumes, err := uc.repo.GetByCandidateId(ctx, req.CandidateId)

	if err != nil {
		return nil, err
	}

	var result []*GetByCandidateIdResponse
	for _, resume := range resumes {
		item := &GetByCandidateIdResponse{
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
	// Convert time.Time to string
	dateOfBirthStr := domainData.DateOfBirth.Format("02.01.2006")

	responseData := ResumeData{
		LastName:        domainData.LastName,
		FirstName:       domainData.FirstName,
		MiddleName:      domainData.MiddleName,
		DateOfBirth:     dateOfBirthStr,
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
