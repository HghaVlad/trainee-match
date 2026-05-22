package application

import "errors"

var (
	ErrActiveAlreadyExists     = errors.New("active application already exists")
	ErrInvalidStatus           = errors.New("invalid status")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrStatusAlreadySet        = errors.New("status already set")
	ErrNotFound                = errors.New("application not found")
	ErrResumeNotPublished      = errors.New("resume must be published")
	ErrResumeAccessDenied      = errors.New("resume access denied")
	ErrAccessDenied            = errors.New("access denied")

	ErrStatusChangeRequiresUserID = errors.New("status change with this actor requires user ID")
	ErrCoverLetterTooLong         = errors.New("cover letter too long")
	ErrStatusChangeCommentTooLong = errors.New("status change comment too long")
)
