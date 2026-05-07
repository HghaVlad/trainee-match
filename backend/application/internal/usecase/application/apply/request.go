package apply

import "github.com/google/uuid"

type Request struct {
	VacancyID   uuid.UUID
	ResumeID    uuid.UUID
	CoverLetter *string
}
