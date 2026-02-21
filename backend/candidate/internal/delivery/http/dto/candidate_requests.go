package dto

import (
	"errors"
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
	return nil
}

type CandidateUpdateRequest struct {
	Phone    *string `json:"phone"`
	Telegram *string `json:"telegram"`
	City     *string `json:"city"`
	Birthday *Date   `json:"birthday"`
}

func (req *CandidateUpdateRequest) Validate() error {
	if req.Phone != nil && *req.Phone == "" {
		return errors.New("phone is required")
	}
	if req.Telegram != nil && *req.Telegram == "" {
		return errors.New("telegram is required")
	}
	if req.City != nil && *req.City == "" {
		return errors.New("city is required")
	}
	return nil
}
