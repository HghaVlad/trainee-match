package dto

import (
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/get_resume"
	"github.com/google/uuid"
)

type ResumeResponse struct {
	ID          uuid.UUID  `json:"id"`
	CandidateID uuid.UUID  `json:"candidate_id"`
	Name        string     `json:"name"`
	Status      int        `json:"status"`
	Data        ResumeData `json:"data"`
}

func UseCaseResponseToDtoResumeResponse(resp get_resume.Response) ResumeResponse {
	educationDTO := make([]Education, len(resp.Data.Education))
	for i, edu := range resp.Data.Education {
		educationDTO[i] = Education{
			Level:          edu.Level,
			University:     edu.University,
			Faculty:        edu.Faculty,
			Specialization: edu.Specialization,
			StartYear:      edu.StartYear,
			EndYear:        edu.EndYear,
			Format:         edu.Format,
		}
	}

	workExpDTO := make([]WorkExperience, len(resp.Data.WorkExperiences))
	for i, exp := range resp.Data.WorkExperiences {
		workExpDTO[i] = WorkExperience{
			Position:         exp.Position,
			Company:          exp.Company,
			Period:           exp.Period,
			Responsibilities: exp.Responsibilities,
		}
	}

	// Parse date string back to DTO Date type
	dtoData := ResumeData{
		LastName:        resp.Data.LastName,
		FirstName:       resp.Data.FirstName,
		MiddleName:      resp.Data.MiddleName,
		DateOfBirth:     TimeToDate(resp.Data.DateOfBirth),
		Email:           resp.Data.Email,
		Phone:           resp.Data.Phone,
		City:            resp.Data.City,
		Citizenship:     resp.Data.Citizenship,
		Education:       educationDTO,
		WorkExperiences: workExpDTO,
		SkillsList:      resp.Data.SkillsList,
		AdditionalInfo:  resp.Data.AdditionalInfo,
		PortfolioLink:   resp.Data.PortfolioLink,
		DesiredFormat:   resp.Data.DesiredFormat,
		EnglishLevel:    resp.Data.EnglishLevel,
	}

	response := ResumeResponse{
		ID:          resp.ID,
		CandidateID: resp.CandidateID,
		Name:        resp.Name,
		Status:      resp.Status,
		Data:        dtoData,
	}
	return response
}

type ShortResumeResponse struct {
	ID          uuid.UUID `json:"id"`
	CandidateId uuid.UUID `json:"candidate_id"`
	Name        string    `json:"name"`
	Status      int       `json:"status"`
}

type SkillResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
