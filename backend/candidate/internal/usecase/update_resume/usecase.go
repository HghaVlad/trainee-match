package update_resume

import (
	"context"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/google/uuid"
	"time"
)

type ResumeRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (domain.Resume, error)
	Update(ctx context.Context, resume *domain.Resume) error
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
	resume, err := uc.resumeRepo.GetById(ctx, req.ID)
	if err != nil {
		return Response{}, err
	}

	if req.Name != nil {
		resume.Name = *req.Name
	}
	if req.Status != nil {
		resume.Status = *req.Status
	}
	if req.Data != nil {
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

		// Convert request data to domain data
		dateOfBirth, err := time.Parse("02.01.2006", req.Data.DateOfBirth)
		if err != nil {
			return Response{}, err
		}

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

		resume.Data = domainData
	}

	err = uc.resumeRepo.Update(ctx, &resume)
	if err != nil {
		return Response{}, err
	}

	return Response{Success: true}, nil
}