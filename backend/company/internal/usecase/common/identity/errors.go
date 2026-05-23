package identity

import "errors"

var (
	ErrHrRoleRequired    = errors.New("hr role is required")
	ErrAdminRoleRequired = errors.New("admin role is required")
	ErrInsufficientRole  = errors.New("insufficient role")
)
