package application

import "errors"

var (
	ErrActiveAlreadyExists     = errors.New("active application already exists")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrStatusAlreadySet        = errors.New("status already set")
	ErrNotFound                = errors.New("application not found")
	ErrVacancyNotPublished     = errors.New("vacancy must be published")
	ErrResumeNotPublished      = errors.New("resume must be published")
	ErrResumeAccessDenied      = errors.New("resume access denied")
	ErrAccessDenied            = errors.New("access denied")
)

var (
	ErrStatusChangeRequiresUserID = errors.New("status change with this actor requires user ID")
)
