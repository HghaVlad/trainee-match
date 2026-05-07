package projection

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type VacancyStatus string

const (
	VacancyStatusPublished VacancyStatus = "published"
	VacancyStatusArchived  VacancyStatus = "archived"
)

type Vacancy struct {
	ID          uuid.UUID
	CompanyID   uuid.UUID
	CompanyName string
	Title       string
	Status      VacancyStatus
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

var (
	ErrVacancyNotFound = errors.New("vacancy not found")
)
