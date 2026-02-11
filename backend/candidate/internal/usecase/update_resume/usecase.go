package update_resume

import (
	"context"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/google/uuid"
)

type ResumeRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (domain.Resume, error)
	Update(ctx context.Context, resume *domain.Resume) error
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
		if req.Data.SkillsList != nil {
			ok, err := uc.skillRepo.AreSkillsExist(ctx, *req.Data.SkillsList)
			if err != nil {
				return Response{}, err
			}
			if !ok {
				return Response{}, domain.ErrSkillNotFound
			}
			resume.Data.SkillsList = *req.Data.SkillsList
		}

		if req.Data.LastName != nil {
			resume.Data.LastName = *req.Data.LastName
		}
		if req.Data.FirstName != nil {
			resume.Data.FirstName = *req.Data.FirstName
		}
		if req.Data.MiddleName != nil {
			resume.Data.MiddleName = *req.Data.MiddleName
		}
		if req.Data.DateOfBirth != nil {
			resume.Data.DateOfBirth = *req.Data.DateOfBirth
		}
		if req.Data.Email != nil {
			resume.Data.Email = *req.Data.Email
		}
		if req.Data.Phone != nil {
			resume.Data.Phone = *req.Data.Phone
		}
		if req.Data.City != nil {
			resume.Data.City = *req.Data.City
		}
		if req.Data.Citizenship != nil {
			resume.Data.Citizenship = *req.Data.Citizenship
		}
		if req.Data.AdditionalInfo != nil {
			resume.Data.AdditionalInfo = *req.Data.AdditionalInfo
		}
		if req.Data.PortfolioLink != nil {
			resume.Data.PortfolioLink = *req.Data.PortfolioLink
		}
		if req.Data.DesiredFormat != nil {
			resume.Data.DesiredFormat = *req.Data.DesiredFormat
		}
		if req.Data.EnglishLevel != nil {
			resume.Data.EnglishLevel = *req.Data.EnglishLevel
		}

		// For slices, we replace the entire slice if provided in the request
		if req.Data.Education != nil {
			resume.Data.Education = make([]domain.Education, len(*req.Data.Education))
			for i, edu := range *req.Data.Education {
				resume.Data.Education[i] = domain.Education{
					Level:          edu.Level,
					University:     edu.University,
					Faculty:        edu.Faculty,
					Specialization: edu.Specialization,
					StartYear:      edu.StartYear,
					EndYear:        edu.EndYear,
					Format:         edu.Format,
				}
			}
		}

		if req.Data.WorkExperiences != nil {
			resume.Data.WorkExperiences = make([]domain.WorkExperience, len(*req.Data.WorkExperiences))
			for i, exp := range *req.Data.WorkExperiences {
				resume.Data.WorkExperiences[i] = domain.WorkExperience{
					Position:         exp.Position,
					Company:          exp.Company,
					Period:           exp.Period,
					Responsibilities: exp.Responsibilities,
				}
			}
		}
	}

	err = uc.resumeRepo.Update(ctx, &resume)
	if err != nil {
		return Response{}, err
	}

	return Response{Success: true}, nil
}
