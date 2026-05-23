package projection

import (
	"time"

	"github.com/google/uuid"
)

type ResumeStatus string

const (
	ResumeStatusDraft     ResumeStatus = "draft"
	ResumeStatusPublished ResumeStatus = "published"
)

type ResumeData struct {
	LastName        string           `json:"last_name"        avro:"last_name"`
	FirstName       string           `json:"first_name"       avro:"first_name"`
	MiddleName      string           `json:"middle_name"      avro:"middle_name"`
	DateOfBirth     time.Time        `json:"date_of_birth"    avro:"date_of_birth"`
	Email           string           `json:"email"            avro:"email"`
	Phone           string           `json:"phone"            avro:"phone"`
	City            string           `json:"city"             avro:"city"`
	Citizenship     string           `json:"citizenship"      avro:"citizenship"`
	Education       []Education      `json:"education"        avro:"education"`
	WorkExperiences []WorkExperience `json:"work_experiences" avro:"work_experiences"`
	SkillsList      []uuid.UUID      `json:"skills_list"      avro:"skills_list"`
	AdditionalInfo  string           `json:"additional_info"  avro:"additional_info"`
	PortfolioLink   string           `json:"portfolio_link"   avro:"portfolio_link"`
	DesiredFormat   string           `json:"desired_format"   avro:"desired_format"`
	EnglishLevel    string           `json:"english_level"    avro:"english_level"`
}

type Education struct {
	Level          string `json:"level"          avro:"level"`
	University     string `json:"university"     avro:"university"`
	Faculty        string `json:"faculty"        avro:"faculty"`
	Specialization string `json:"specialization" avro:"specialization"`
	StartYear      int    `json:"start_year"     avro:"start_year"`
	EndYear        int    `json:"end_year"       avro:"end_year"`
	Format         string `json:"format"         avro:"format"`
}

type WorkExperience struct {
	Position         string `json:"position"         avro:"position"`
	Company          string `json:"company"          avro:"company"`
	Period           string `json:"period"           avro:"period"`
	Responsibilities string `json:"responsibilities" avro:"responsibilities"`
}
