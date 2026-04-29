package domain

import (
	"time"

	"github.com/google/uuid"
)

// ApplicationActor mirrors application_actor_enum.
type ApplicationActor string

const (
	ApplicationActorCandidate ApplicationActor = "candidate"
	ApplicationActorHR        ApplicationActor = "hr"
	ApplicationActorSystem    ApplicationActor = "system"
)

// ApplicationStatusHistory mirrors application_status_history.
type ApplicationStatusHistory struct {
	ID uuid.UUID

	ApplicationID uuid.UUID

	Status ApplicationStatus

	ChangedByUserID uuid.UUID
	ChangedByRole   ApplicationActor

	Comment string

	CreatedAt time.Time
}
