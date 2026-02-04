package dto

import (
	"errors"
	"github.com/google/uuid"
	"regexp"
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

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(r.Email) {
		return errors.New("invalid email format")
	}

	// Validate phone format
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	if !phoneRegex.MatchString(r.Phone) {
		return errors.New("invalid phone number format")
	}

	return nil
}
