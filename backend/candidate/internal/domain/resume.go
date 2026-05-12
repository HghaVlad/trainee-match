package domain

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrResumeNotFound             = errors.New("resume not found")
	ErrInvalidResumeName          = errors.New("name is required")
	ErrInvalidResumeStatus        = errors.New("invalid status")
	ErrInvalidName                = errors.New("first name and last name are required")
	ErrInvalidEmailFormat         = errors.New("invalid email format")
	ErrDateOfBirthInFuture        = errors.New("date of birth cannot be in the future")
	ErrInvalidCitizenship         = errors.New("citizenship is required")
	ErrInvalidEducationEntry      = errors.New("invalid education entry")
	ErrInvalidWorkExperienceEntry = errors.New("invalid work experience entry")
	ErrInvalidPortfolioLink       = errors.New("invalid portfolio link")
)

type Resume struct {
	ID          uuid.UUID
	CandidateId uuid.UUID
	Name        string
	Status      int
	Data        ResumeData
}

type Education struct {
	Level          string `avro:"level"`
	University     string `avro:"university"`
	Faculty        string `avro:"faculty"`
	Specialization string `avro:"specialization"`
	StartYear      int    `avro:"start_year"`
	EndYear        int    `avro:"end_year"`
	Format         string `avro:"format"`
}

type WorkExperience struct {
	Position         string `avro:"position"`
	Company          string `avro:"company"`
	Period           string `avro:"period"`
	Responsibilities string `avro:"responsibilities"`
}

type ResumeData struct {
	LastName        string           `avro:"last_name"`
	FirstName       string           `avro:"first_name"`
	MiddleName      string           `avro:"middle_name"`
	DateOfBirth     time.Time        `avro:"date_of_birth"`
	Email           string           `avro:"email"`
	Phone           string           `avro:"phone"`
	City            string           `avro:"city"`
	Citizenship     string           `avro:"citizenship"`
	Education       []Education      `avro:"education"`
	WorkExperiences []WorkExperience `avro:"work_experiences"`
	SkillsList      []uuid.UUID      `avro:"skills_list"`
	AdditionalInfo  string           `avro:"additional_info"`
	PortfolioLink   string           `avro:"portfolio_link"`
	DesiredFormat   string           `avro:"desired_format"`
	EnglishLevel    string           `avro:"english_level"`
}

// Validate for Education entry
func (e Education) Validate() error {
	if strings.TrimSpace(e.Level) == "" || strings.TrimSpace(e.University) == "" {
		return ErrInvalidEducationEntry
	}
	if e.StartYear < 1900 || e.EndYear < 1900 || e.StartYear > e.EndYear {
		return ErrInvalidEducationEntry
	}
	return nil
}

// Validate for WorkExperience entry
func (w WorkExperience) Validate() error {
	if strings.TrimSpace(w.Position) == "" || strings.TrimSpace(w.Company) == "" {
		return ErrInvalidWorkExperienceEntry
	}
	return nil
}

// Validate checks whole resume for domain-level business rules.
func (r Resume) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return ErrInvalidResumeName
	}

	if err := r.Data.Validate(); err != nil {
		return err
	}

	for _, id := range r.Data.SkillsList {
		if id == uuid.Nil {
			return errors.New("invalid skill id")
		}
	}
	return nil
}

// Validate checks fields inside ResumeData for correctness.
func (d ResumeData) Validate() error {
	if strings.TrimSpace(d.FirstName) == "" || strings.TrimSpace(d.LastName) == "" {
		return ErrInvalidName
	}

	if !d.DateOfBirth.IsZero() && d.DateOfBirth.After(time.Now()) {
		return ErrDateOfBirthInFuture
	}

	if strings.TrimSpace(d.Email) == "" {
		return ErrInvalidEmailFormat
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(d.Email) {
		return ErrInvalidEmailFormat
	}

	if strings.TrimSpace(d.Phone) == "" {
		return ErrInvalidPhoneFormat
	}
	phoneRegex := regexp.MustCompile(`^[+\d()\-\s]{7,25}$`)
	if !phoneRegex.MatchString(d.Phone) {
		return ErrInvalidPhoneFormat
	}

	if strings.TrimSpace(d.City) == "" {
		return ErrInvalidCityFormat
	}
	if strings.TrimSpace(d.Citizenship) == "" {
		return ErrInvalidCitizenship
	}

	for _, edu := range d.Education {
		if err := edu.Validate(); err != nil {
			return err
		}
	}

	for _, we := range d.WorkExperiences {
		if err := we.Validate(); err != nil {
			return err
		}
	}

	// Portfolio link if present
	if strings.TrimSpace(d.PortfolioLink) != "" {
		urlRegex := regexp.MustCompile(`^https?://`)
		if !urlRegex.MatchString(d.PortfolioLink) {
			return ErrInvalidPortfolioLink
		}
	}
	return nil
}
