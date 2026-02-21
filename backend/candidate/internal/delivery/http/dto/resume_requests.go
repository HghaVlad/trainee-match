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
	if err := req.Data.Validate(); err != nil {
		return err
	}

	return nil
}

type UpdateResumeRequest struct {
	ID     *uuid.UUID       `json:"id"`
	Name   *string          `json:"name"`
	Status *int             `json:"status"`
	Data   *PatchResumeData `json:"data"`
}

func (req *UpdateResumeRequest) Validate() error {
	if req.Data != nil {
		if err := req.Data.Validate(); err != nil {
			return err
		}
	}

	return nil
}
