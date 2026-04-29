package domain

import "strings"

type ResumeStatus string

const (
	Draft     ResumeStatus = "draft"
	Published ResumeStatus = "published"
)

func Parse(status string) (int, error) {
	normalized := strings.ToLower(strings.TrimSpace(status))
	switch normalized {
	case string(Draft):
		return 0, nil
	case string(Published):
		return 1, nil
	default:
		return -1, ErrInvalidResumeStatus
	}
}

func Format(status int) (ResumeStatus, error) {
	switch status {
	case 0:
		return Draft, nil
	case 1:
		return Published, nil
	default:
		return "", ErrInvalidResumeStatus
	}
}
