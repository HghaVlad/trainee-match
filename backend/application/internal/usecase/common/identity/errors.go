package identity

import "errors"

var (
	ErrHrRoleRequired        = errors.New("hr role is required")
	ErrCandidateRoleRequired = errors.New("candidate role is required")
)
