package identity

import "errors"

var ErrHrRoleRequired = errors.New("hr role is required")
var ErrInsufficientRole = errors.New("insufficient role")
