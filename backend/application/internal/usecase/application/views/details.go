package views

import (
	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
)

type AllowedAction string

const (
	AllowedActionWithdraw AllowedAction = "withdraw"
)

type Details struct {
	Application    *application.Application
	VacProj        *projection.Vacancy
	Snapshot       *application.Snapshot
	StatusHistory  []application.StatusChange
	AllowedActions []AllowedAction
}
