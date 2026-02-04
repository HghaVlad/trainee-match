package create_resume

import (
	"context"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/google/uuid"
	"time"
)

type ResumeRepo interface {
	Create(ctx context.Context, resume *domain.Resume) (uuid.UUID, error)
}

type SkillRepo interface {
	CheckExistsBatch(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]bool, error)
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
	if len(req.Data.SkillsList) > 0 {
		existingSkills, err := uc.skillRepo.CheckExistsBatch(ctx, req.Data.SkillsList)
		if err != nil {
			return Response{}, err
		}

		// Check if any skill doesn't exist
		for _, skillId := range req.Data.SkillsList {
			if !existingSkills[skillId] {
				return Response{}, domain.ErrSkillNotFound
			}
		}
	}

	// Convert date string to time.Time
	dateOfBirth, err := time.Parse("02.01.2006", req.Data.DateOfBirth)
	if err != nil {
		return Response{}, err
	}

	// Convert request data to domain model
	domainData := domain.ResumeData{
		LastName:        req.Data.LastName,
		FirstName:       req.Data.FirstName,
		MiddleName:      req.Data.MiddleName,
		DateOfBirth:     dateOfBirth,
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
