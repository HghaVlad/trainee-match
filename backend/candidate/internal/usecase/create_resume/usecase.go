package create_resume

import (
	"context"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/google/uuid"
)

type ResumeRepo interface {
	Create(ctx context.Context, resume *domain.Resume) (uuid.UUID, error)
}

type SkillRepo interface {
	AreSkillsExist(ctx context.Context, ids []uuid.UUID) (bool, error)
}

type UseCase struct {
	resumeRepo ResumeRepo
	skillRepo  SkillRepo
}

func New(resumeRepo ResumeRepo, skillRepo SkillRepo) *UseCase {
	return &UseCase{
		resumeRepo: resumeRepo,
		skillRepo:  skillRepo,
	}
}

func (uc *UseCase) Execute(ctx context.Context, req Request) (Response, error) {
	// Check if all skills exist
	ok, err := uc.skillRepo.AreSkillsExist(ctx, req.Data.SkillsList)
	if err != nil {
		return Response{}, err
	}
	if !ok {
		return Response{}, domain.ErrSkillNotFound
	}

	// Convert request data to domain model
	domainData := domain.ResumeData{
		LastName:        req.Data.LastName,
		FirstName:       req.Data.FirstName,
		MiddleName:      req.Data.MiddleName,
		DateOfBirth:     req.Data.DateOfBirth,
		Email:           req.Data.Email,
		Phone:           req.Data.Phone,
		City:            req.Data.City,
		Citizenship:     req.Data.Citizenship,
		Education:       make([]domain.Education, len(req.Data.Education)),
		WorkExperiences: make([]domain.WorkExperience, len(req.Data.WorkExperiences)),
		SkillsList:      req.Data.SkillsList,
		AdditionalInfo:  req.Data.AdditionalInfo,
		PortfolioLink:   req.Data.PortfolioLink,
		DesiredFormat:   req.Data.DesiredFormat,
		EnglishLevel:    req.Data.EnglishLevel,
	}

	for i, edu := range req.Data.Education {
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

	for i, exp := range req.Data.WorkExperiences {
		domainData.WorkExperiences[i] = domain.WorkExperience{
			Position:         exp.Position,
			Company:          exp.Company,
			Period:           exp.Period,
			Responsibilities: exp.Responsibilities,
		}
	}

	resume := &domain.Resume{
		CandidateId: req.CandidateId,
		Name:        req.Name,
		Status:      req.Status,
		Data:        domainData,
	}

	id, err := uc.resumeRepo.Create(ctx, resume)
	if err != nil {
		return Response{}, err
	}

	return Response{ID: id}, nil
}
