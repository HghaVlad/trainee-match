package update_resume

import (
	"time"

	"github.com/google/uuid"
)

type Request struct {
	ID     uuid.UUID   `json:"id"`
	UserId uuid.UUID   `json:"userId"`
	Name   *string     `json:"name"`
	Status *int        `json:"status"`
	Data   *ResumeData `json:"data"`
}

type ResumeData struct {
	LastName        *string           `json:"last_name,omitempty"`
	FirstName       *string           `json:"first_name,omitempty"`
	MiddleName      *string           `json:"middle_name,omitempty"`
	DateOfBirth     *time.Time        `json:"date_of_birth,omitempty"`
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
