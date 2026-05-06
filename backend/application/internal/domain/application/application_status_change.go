package application

import (
	"time"

	"github.com/google/uuid"
)

type Actor string

const (
	ActorCandidate Actor = "candidate"
	ActorHR        Actor = "hr"
	ActorSystem    Actor = "system"
)

type StatusChange struct {
	ID              uuid.UUID
	ApplicationID   uuid.UUID
	Status          Status
	ChangedByUserID *uuid.UUID
	ChangedByRole   Actor
	Comment         *string
	CreatedAt       time.Time
}

func NewStatusChange(
	appID uuid.UUID,
	status Status,
	userID *uuid.UUID,
	actor Actor,
	comment *string,
	createdAt time.Time,
) (*StatusChange, error) {
	if actor != ActorSystem && userID == nil {
		return nil, ErrStatusChangeRequiresUserID
	}

	return &StatusChange{
		ID:              uuid.New(),
		ApplicationID:   appID,
		Status:          status,
		ChangedByUserID: userID,
		ChangedByRole:   actor,
		Comment:         comment,
		CreatedAt:       createdAt,
	}, nil
}
