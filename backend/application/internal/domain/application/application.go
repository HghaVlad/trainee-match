package application

import (
	"time"

	"github.com/google/uuid"
)

const MaxCoverLetterLength = 2000

type Application struct {
	ID          uuid.UUID
	ResumeID    uuid.UUID
	CandidateID uuid.UUID
	VacancyID   uuid.UUID
	CompanyID   uuid.UUID
	SnapshotID  uuid.UUID
	Status      Status
	CoverLetter *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewSubmitted(
	resumeID, candID, vacID, compID, snapID uuid.UUID,
	coverLetter *string,
	when time.Time,
) (*Application, error) {
	if coverLetter != nil && len([]rune(*coverLetter)) > MaxCoverLetterLength {
		return nil, ErrCoverLetterTooLong
	}

	return &Application{
		ID:          uuid.New(),
		ResumeID:    resumeID,
		CandidateID: candID,
		VacancyID:   vacID,
		CompanyID:   compID,
		SnapshotID:  snapID,
		Status:      StatusSubmitted,
		CoverLetter: coverLetter,
		CreatedAt:   when,
		UpdatedAt:   when,
	}, nil
}

func (a *Application) Withdraw(when time.Time) error {
	return a.ChangeStatus(StatusWithdrawn, ActorCandidate, when)
}

func (a *Application) ChangeStatus(next Status, actor Actor, when time.Time) error {
	if a.Status == next {
		return ErrStatusAlreadySet
	}

	if !a.Status.CanTransitionTo(next, actor) {
		return ErrInvalidStatusTransition
	}

	a.Status = next
	a.UpdatedAt = when

	return nil
}
