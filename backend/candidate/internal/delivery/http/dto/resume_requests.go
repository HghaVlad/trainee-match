package dto

import (
	"errors"
	"github.com/google/uuid"
)

type CreateResumeRequest struct {
	Name   string     `json:"name"`
	Status int        `json:"status"`
	Data   ResumeData `json:"data"`
}

func (req *CreateResumeRequest) Validate() error {
	if req.Name == "" {
		return errors.New("name is required")
	}

	if req.Status != 0 && req.Status != 1 { // Assuming 0=Draft, 1=Published
		return errors.New("invalid status")
	}
	if err := req.Data.Validate(); err != nil {
		return err
	}

	return nil
}

type UpdateResumeRequest struct {
	ID     *uuid.UUID  `json:"id"`
	Name   *string     `json:"name"`
	Status *int        `json:"status"`
	Data   *ResumeData `json:"data"`
}

func (req *UpdateResumeRequest) Validate() error {
	if req.ID == nil {
		return errors.New("id is required")
	}

	if req.Status != nil && (*req.Status != 0 && *req.Status != 1) { // Assuming 0=Draft, 1=Published
		return errors.New("invalid status")
	}

	if req.Data != nil {
		if err := req.Data.Validate(); err != nil {
			return err
		}
	}

	return nil
}
