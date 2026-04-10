package identity

import (
	"github.com/google/uuid"
)

type Identity struct {
	UserID uuid.UUID
	Role   GlobalRole
}

type GlobalRole string

const (
	RoleHR        GlobalRole = "Company"
	RoleCandidate GlobalRole = "Candidate"
	RoleAdmin     GlobalRole = "Admin"
)
