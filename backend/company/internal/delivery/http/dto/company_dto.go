package dto

import (
	"time"

	"github.com/google/uuid"
)

type CompanyResponse struct {
	ID               uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name             string    `json:"name" example:"Google Inc."`
	Description      *string   `json:"description,omitempty" example:"We make the world a better place"`
	Website          *string   `json:"website,omitempty" example:"https://www.google.com"`
	LogoURL          *string   `json:"logoURL,omitempty" example:"http://domain/minio/6icinimmck...mksk"`
	OpenVacanciesCnt int       `json:"openVacanciesCount" example:"13"`
	CreatedAt        time.Time `json:"createdAt" example:"2020-04-08T21:00:00Z"`
	UpdatedAt        time.Time `json:"updatedAt" example:"2020-04-08T21:00:00Z"`
}

type CompanyListItemResponse struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	OpenVacanciesCnt int       `json:"openVacanciesCount"`
	LogoURL          *string   `json:"logoUrl,omitempty"`
}

type CompanyListResponse struct {
	Companies  []CompanyListItemResponse `json:"companies"`
	NextCursor *string                   `json:"nextCursor,omitempty"`
}

type CompanyCreateRequest struct {
	Name        string  `json:"name" example:"Google Inc."`
	Description *string `json:"description,omitempty" example:"We make the world a better place"`
	Website     *string `json:"website,omitempty" example:"https://www.google.com"`
}

type CompanyCreatedResponse struct {
	ID uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// CompanyUpdateRequest if Description or Website not provided - don't change
type CompanyUpdateRequest struct {
	Name        *string `json:"name,omitempty" example:"Google LLC"`
	Description *string `json:"description" example:"New description"`
	Website     *string `json:"website" example:"https://google.com"`
}

type CompanyAddHrRequest struct {
	UserID uuid.UUID `json:"userID" example:"550e8400-e29b-41d4-a716-446655440000"`
	Role   string    `json:"role" enums:"recruiter,admin" example:"recruiter"`
}
