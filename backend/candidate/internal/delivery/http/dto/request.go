package dto

import (
	"errors"
	"regexp"
)

type CandidateCreateRequest struct {
	Phone    string `json:"phone"`
	Telegram string `json:"telegram"`
	City     string `json:"city"`
	Birthday Date   `json:"birthday"`
}

func (req *CandidateCreateRequest) Validate() error {
	if req.Phone == "" {
		return errors.New("phone is required")
	}
	if req.Telegram == "" {
		return errors.New("telegram is required")
	}
	if req.City == "" {
		return errors.New("city is required")
	}

	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	if !phoneRegex.MatchString(req.Phone) {
		return errors.New("invalid phone number format")
	}

	telegramRegex := regexp.MustCompile(`^@[\w]{3,32}$`)
	if !telegramRegex.MatchString(req.Telegram) {
		return errors.New("telegram username must start with @ and contain only alphanumeric characters and underscores, 3-32 characters long")
	}

	return nil
}

type CandidateUpdateRequest struct {
	Phone    *string `json:"phone"`
	Telegram *string `json:"telegram"`
	City     *string `json:"city"`
	Birthday *Date   `json:"birthday"`
}

func (req *CandidateUpdateRequest) Validate() error {
	if req.Phone != nil {
		phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
		if !phoneRegex.MatchString(*req.Phone) {
			return errors.New("invalid phone number format")
		}
	}
	if req.Telegram != nil {
		telegramRegex := regexp.MustCompile(`^@[\w]{3,32}$`)
		if !telegramRegex.MatchString(*req.Telegram) {
			return errors.New("telegram username must start with @ and contain only alphanumeric characters and underscores, 3-32 characters long")
		}
	}

	if req.City != nil {
		if *req.City == "" {
			return errors.New("city is required")
		}
	}

	return nil
}
