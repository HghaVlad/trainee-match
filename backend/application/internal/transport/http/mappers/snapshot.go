package mappers

import (
	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
	"github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/oapi"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
)

func snapshotToHTTP(appSnap views.ApplicationSnapshot) oapi.ApplicationSnapshot {
	return oapi.ApplicationSnapshot{
		Email:      emailToOAPI(&appSnap.Email),
		FullName:   appSnap.FullName,
		Telegram:   appSnap.Telegram,
		CreatedAt:  appSnap.CreatedAt,
		ResumeData: resumeDataToHTTP(appSnap.ResumeData),
	}
}

func resumeDataToHTTP(proj projection.ResumeData) oapi.ResumeData {
	education := make([]oapi.Education, 0, len(proj.Education))
	for _, e := range proj.Education {
		education = append(education, oapi.Education{
			Level:          e.Level,
			University:     e.University,
			Faculty:        e.Faculty,
			Specialization: e.Specialization,
			StartYear:      e.StartYear,
			EndYear:        e.EndYear,
			Format:         e.Format,
		})
	}

	workExperiences := make([]oapi.WorkExperience, 0, len(proj.WorkExperiences))
	for _, we := range proj.WorkExperiences {
		workExperiences = append(workExperiences, oapi.WorkExperience{
			Position:         we.Position,
			Company:          we.Company,
			Period:           we.Period,
			Responsibilities: we.Responsibilities,
		})
	}

	skills := make([]uuid.UUID, 0, len(proj.SkillsList))
	skills = append(skills, proj.SkillsList...)

	var additionalInfo *string
	if proj.AdditionalInfo != "" {
		additionalInfo = &proj.AdditionalInfo
	}

	var portfolioLink *string
	if proj.PortfolioLink != "" {
		portfolioLink = &proj.PortfolioLink
	}

	return oapi.ResumeData{
		LastName:        proj.LastName,
		FirstName:       proj.FirstName,
		MiddleName:      proj.MiddleName,
		DateOfBirth:     openapitypes.Date{Time: proj.DateOfBirth},
		Email:           openapitypes.Email(proj.Email),
		Phone:           proj.Phone,
		City:            proj.City,
		Citizenship:     proj.Citizenship,
		Education:       education,
		WorkExperiences: workExperiences,
		SkillsList:      skills,
		AdditionalInfo:  additionalInfo,
		PortfolioLink:   portfolioLink,
		DesiredFormat:   proj.DesiredFormat,
		EnglishLevel:    proj.EnglishLevel,
	}
}
