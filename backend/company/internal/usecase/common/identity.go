package uc_common

import "github.com/google/uuid"

type Identity struct {
	UserID uuid.UUID
	Role   GlobalRole
}

type GlobalRole string

const (
	RoleHR        GlobalRole = "HR"
	RoleCandidate GlobalRole = "Candidate"
	RoleAdmin     GlobalRole = "Admin"
)
