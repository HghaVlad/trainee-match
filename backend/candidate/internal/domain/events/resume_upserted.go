package events

import (
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

type ResumeUpserted struct {
	EventID     uuid.UUID           `avro:"event_id"`
	ResumeID    uuid.UUID           `avro:"resume_id"`
	OccurredAt  time.Time           `avro:"occurred_at"`
	CandidateID uuid.UUID           `avro:"candidate_id"`
	Name        string              `avro:"name"`
	Data        domain.ResumeData   `avro:"data"`
	Status      domain.ResumeStatus `avro:"status"`
	CreatedAt   *time.Time          `avro:"created_at"`
	UpdatedAt   *time.Time          `avro:"updated_at"`
}

func NewResumeUpserted(r domain.Resume) ResumeUpserted {
	resumeStatus, _ := domain.Format(r.Status)
	return ResumeUpserted{
		EventID:     uuid.New(),
		ResumeID:    r.ID,
		OccurredAt:  time.Now().UTC(),
		CandidateID: r.CandidateId,
		Name:        r.Name,
		Status:      resumeStatus,
		CreatedAt:   nil,
		UpdatedAt:   nil,
		Data: domain.ResumeData{
			LastName:        r.Data.LastName,
			FirstName:       r.Data.FirstName,
			MiddleName:      r.Data.MiddleName,
			DateOfBirth:     r.Data.DateOfBirth,
			Email:           r.Data.Email,
			Phone:           r.Data.Phone,
			City:            r.Data.City,
			Citizenship:     r.Data.Citizenship,
			Education:       r.Data.Education,
			WorkExperiences: r.Data.WorkExperiences,
			SkillsList:      r.Data.SkillsList,
			AdditionalInfo:  r.Data.AdditionalInfo,
			PortfolioLink:   r.Data.PortfolioLink,
			DesiredFormat:   r.Data.DesiredFormat,
			EnglishLevel:    r.Data.EnglishLevel,
		},
	}
}
