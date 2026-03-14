package dto

import (
	"errors"
	"time"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/create_resume"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/update_resume"
	"github.com/google/uuid"
)

type CreateResumeRequest struct {
	Name   string     `json:"name"`
	Status int        `json:"status"`
	Data   ResumeData `json:"data"`
}

func (req *CreateResumeRequest) Validate() error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if err := req.Data.Validate(); err != nil {
		return err
	}

	return nil
}

func (req *CreateResumeRequest) ToUseCaseRequest() create_resume.Request {
	useCaseReq := create_resume.Request{
		Name:   req.Name,
		Status: req.Status,
		Data: create_resume.ResumeData{
			LastName:        req.Data.LastName,
			FirstName:       req.Data.FirstName,
			MiddleName:      req.Data.MiddleName,
			DateOfBirth:     DateToTime(req.Data.DateOfBirth),
			Email:           req.Data.Email,
			Phone:           req.Data.Phone,
			City:            req.Data.City,
			Citizenship:     req.Data.Citizenship,
			Education:       make([]create_resume.Education, len(req.Data.Education)),
			WorkExperiences: make([]create_resume.WorkExperience, len(req.Data.WorkExperiences)),
			SkillsList:      req.Data.SkillsList,
			AdditionalInfo:  req.Data.AdditionalInfo,
			PortfolioLink:   req.Data.PortfolioLink,
			DesiredFormat:   req.Data.DesiredFormat,
			EnglishLevel:    req.Data.EnglishLevel,
		},
	}

	for i, edu := range req.Data.Education {
		useCaseReq.Data.Education[i] = create_resume.Education{
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
		useCaseReq.Data.WorkExperiences[i] = create_resume.WorkExperience{
			Position:         exp.Position,
			Company:          exp.Company,
			Period:           exp.Period,
			Responsibilities: exp.Responsibilities,
		}
	}
	return useCaseReq
}

type UpdateResumeRequest struct {
	ID     *uuid.UUID       `json:"id"`
	Name   *string          `json:"name"`
	Status *int             `json:"status"`
	Data   *PatchResumeData `json:"data"`
}

func (req *UpdateResumeRequest) Validate() error {
	if req.Data != nil {
		if err := req.Data.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (req *UpdateResumeRequest) ToUseCaseRequest() update_resume.Request {
	useCaseReq := update_resume.Request{
		Name:   req.Name,
		Status: req.Status,
	}

	if req.Data != nil {
		var dateOfBirthStr time.Time
		if req.Data.DateOfBirth != nil {
			dateOfBirthStr = DateToTime(*req.Data.DateOfBirth)
		}

		useCaseReq.Data = &update_resume.ResumeData{
			LastName:       req.Data.LastName,
			FirstName:      req.Data.FirstName,
			MiddleName:     req.Data.MiddleName,
			DateOfBirth:    &dateOfBirthStr,
			Email:          req.Data.Email,
			Phone:          req.Data.Phone,
			City:           req.Data.City,
			Citizenship:    req.Data.Citizenship,
			AdditionalInfo: req.Data.AdditionalInfo,
			PortfolioLink:  req.Data.PortfolioLink,
			DesiredFormat:  req.Data.DesiredFormat,
			EnglishLevel:   req.Data.EnglishLevel,
		}

		// Handle Education if provided
		if req.Data.Education != nil {
			educationSlice := make([]update_resume.Education, len(*req.Data.Education))
			for i, edu := range *req.Data.Education {
				educationSlice[i] = update_resume.Education{
					Level:          edu.Level,
					University:     edu.University,
					Faculty:        edu.Faculty,
					Specialization: edu.Specialization,
					StartYear:      edu.StartYear,
					EndYear:        edu.EndYear,
					Format:         edu.Format,
				}
			}
			useCaseReq.Data.Education = &educationSlice
		}

		// Handle WorkExperiences if provided
		if req.Data.WorkExperiences != nil {
			workExpSlice := make([]update_resume.WorkExperience, len(*req.Data.WorkExperiences))
			for i, exp := range *req.Data.WorkExperiences {
				workExpSlice[i] = update_resume.WorkExperience{
					Position:         exp.Position,
					Company:          exp.Company,
					Period:           exp.Period,
					Responsibilities: exp.Responsibilities,
				}
			}
			useCaseReq.Data.WorkExperiences = &workExpSlice
		}

		// Handle SkillsList if provided
		if req.Data.SkillsList != nil {
			skillsListCopy := make([]uuid.UUID, len(*req.Data.SkillsList))
			copy(skillsListCopy, *req.Data.SkillsList)
			useCaseReq.Data.SkillsList = &skillsListCopy
		}
	}
	return useCaseReq
}
