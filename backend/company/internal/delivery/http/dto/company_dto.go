package dto

import (
	"time"

	"github.com/google/uuid"
)

type CompanyResponse struct {
	ID          uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string    `json:"name" example:"Google Inc."`
	Description *string   `json:"description,omitempty" example:"We make the world a better place"`
	Website     *string   `json:"website,omitempty" example:"https://www.google.com"`
	OwnerId     uuid.UUID `json:"ownerId" example:"550e8400-e29b-41d4-a716-446655440000"`
	LogoURL     *string   `json:"logoURL,omitempty" example:"http://domain/minio/6icinimmck...mksk"`
	CreatedAt   time.Time `json:"createdAt" example:"2020-04-08T21:00:00Z"`
	UpdatedAt   time.Time `json:"updatedAt" example:"2020-04-08T21:00:00Z"`
}
