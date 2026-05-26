package get_resume

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

type ResumeRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (domain.Resume, error)
	GetByCandidateId(ctx context.Context, candidateId uuid.UUID, page, size int) ([]domain.Resume, error)
}

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

func (uc *UseCase) GetById(ctx context.Context, resumeId, UserId uuid.UUID) (*Response, error) {
	candidate, err := uc.candidateRepo.GetByUserID(ctx, UserId)
	if err != nil {
		return nil, err
	}
	resume, err := uc.resumeRepo.GetById(ctx, resumeId)

	if err != nil {
		return nil, err
	}

	if candidate.ID != resume.CandidateId || resume.ModerationStatus == domain.ModerationStatusHidden {
		return nil, domain.ErrResumeNotFound
	}

	status, err := domain.Format(resume.Status)
	if err != nil {
		return nil, err
	}

	response := &Response{
		ID:               resume.ID,
		CandidateID:      resume.CandidateId,
		Name:             resume.Name,
		Status:           string(status),
		ModerationStatus: string(resume.ModerationStatus),
		Data:             convertDomainDataToResponseData(resume.Data),
	}

	return response, nil
}

func (uc *UseCase) GetByCandidateId(ctx context.Context, UserId uuid.UUID, page, size int) ([]*ShortResponse, error) {
	candidate, err := uc.candidateRepo.GetByUserID(ctx, UserId)
	if err != nil {
		return nil, err
	}
	resumes, err := uc.resumeRepo.GetByCandidateId(ctx, candidate.ID, page, size)

	if err != nil {
		return nil, err
	}

	var result []*ShortResponse
	for _, resume := range resumes {
		status, err := domain.Format(resume.Status)
		if err != nil {
			return nil, err
		}
		if resume.ModerationStatus == domain.ModerationStatusHidden {
			continue
		}
		item := &ShortResponse{
			ID:               resume.ID,
			CandidateId:      resume.CandidateId,
			Name:             resume.Name,
			Status:           string(status),
			ModerationStatus: string(resume.ModerationStatus),
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
