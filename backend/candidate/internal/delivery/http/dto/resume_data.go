package dto

import (
	"errors"

	"github.com/google/uuid"
)

type ResumeData struct {
	LastName        string           `json:"last_name"`
	FirstName       string           `json:"first_name"`
	MiddleName      string           `json:"middle_name"`
	DateOfBirth     Date             `json:"date_of_birth"`
	Email           string           `json:"email"`
	Phone           string           `json:"phone"`
	City            string           `json:"city"`
	Citizenship     string           `json:"citizenship"`
	Education       []Education      `json:"education"`
	WorkExperiences []WorkExperience `json:"work_experiences"`
	SkillsList      []uuid.UUID      `json:"skills_list"` // Store skill IDs
	AdditionalInfo  string           `json:"additional_info"`
	PortfolioLink   string           `json:"portfolio_link"`
	DesiredFormat   string           `json:"desired_format"`
	EnglishLevel    string           `json:"english_level"`
}

type Education struct {
	Level          string `json:"level"`
	University     string `json:"university"`
	Faculty        string `json:"faculty"`
	Specialization string `json:"specialization"`
	StartYear      int    `json:"start_year"`
	EndYear        int    `json:"end_year"`
	Format         string `json:"format"`
}

type WorkExperience struct {
	Position         string `json:"position"`
	Company          string `json:"company"`
	Period           string `json:"period"`
	Responsibilities string `json:"responsibilities"`
}

func (r *ResumeData) Validate() error {
	if r.LastName == "" {
		return errors.New("last name is required")
	}
	if r.FirstName == "" {
		return errors.New("first name is required")
	}

	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Phone == "" {
		return errors.New("phone is required")
	}
	if r.City == "" {
		return errors.New("city is required")
	}
	if r.DateOfBirth == (Date{}) {
		return errors.New("date of birth is required")
	}
	if r.Citizenship == "" {
		return errors.New("citizenship is required")
	}
	if r.DesiredFormat == "" {
		return errors.New("desired format is required")
	}
	if r.EnglishLevel == "" {
		return errors.New("english level is required")
	}

	return nil
}

type PatchResumeData struct {
	LastName        *string           `json:"last_name,omitempty"`
	FirstName       *string           `json:"first_name,omitempty"`
	MiddleName      *string           `json:"middle_name,omitempty"`
	DateOfBirth     *Date             `json:"date_of_birth,omitempty"`
	Email           *string           `json:"email,omitempty"`
	Phone           *string           `json:"phone,omitempty"`
	City            *string           `json:"city,omitempty"`
	Citizenship     *string           `json:"citizenship,omitempty"`
	Education       *[]Education      `json:"education,omitempty"`
	WorkExperiences *[]WorkExperience `json:"work_experiences,omitempty"`
	SkillsList      *[]uuid.UUID      `json:"skills_list,omitempty"`
	AdditionalInfo  *string           `json:"additional_info,omitempty"`
	PortfolioLink   *string           `json:"portfolio_link,omitempty"`
	DesiredFormat   *string           `json:"desired_format,omitempty"`
	EnglishLevel    *string           `json:"english_level,omitempty"`
}

func (r *PatchResumeData) Validate() error {
	if r.LastName != nil && *r.LastName == "" {
		return errors.New("last name cannot be empty")
	}
	if r.FirstName != nil && *r.FirstName == "" {
		return errors.New("first name cannot be empty")
	}
	if r.City != nil && *r.City == "" {
		return errors.New("city cannot be empty")
	}

	return nil
}
