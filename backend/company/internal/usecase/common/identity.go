package uc_common

import "github.com/google/uuid"

type Identity struct {
	UserID uuid.UUID
	Role   GlobalRole
}

type GlobalRole string

const (
	RoleHR        GlobalRole = "Company" // TODO: ask vlad to change it to HR or smth and to place it normally
	RoleCandidate GlobalRole = "Candidate"
	RoleAdmin     GlobalRole = "Admin"
)
