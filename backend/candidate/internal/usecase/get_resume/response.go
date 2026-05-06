package get_resume

import (
	"time"

	"github.com/google/uuid"
)

type Response struct {
	ID          uuid.UUID  `json:"id"`
	CandidateID uuid.UUID  `json:"candidate_id"`
	Name        string     `json:"name"`
	Status      string     `json:"status"`
	Data        ResumeData `json:"data"`
}

type ShortResponse struct {
	ID          uuid.UUID `json:"id"`
	CandidateId uuid.UUID `json:"candidate_id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
}

type ResumeData struct {
	LastName        string           `json:"last_name"`
	FirstName       string           `json:"first_name"`
	MiddleName      string           `json:"middle_name"`
	DateOfBirth     time.Time        `json:"date_of_birth"` // Using string to represent date
	Email           string           `json:"email"`
	Phone           string           `json:"phone"`
	City            string           `json:"city"`
	Citizenship     string           `json:"citizenship"`
	Education       []Education      `json:"education"`
	WorkExperiences []WorkExperience `json:"work_experiences"`
	SkillsList      []uuid.UUID      `json:"skills_list"`
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
