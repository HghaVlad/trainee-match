package remove_resume

import "github.com/google/uuid"

type Request struct {
	ResumeId uuid.UUID
	UserId   uuid.UUID
}
