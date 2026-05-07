package listcandidatesummary

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
)

type Request struct {
	Statuses  []application.Status
	CompanyID *uuid.UUID
	Cursor    string
	Limit     int
	Order     Order
}

type Order string

const (
	OrderCreatedAtDesc = "createdAtDesc"
	OrderUpdatedAtDesc = "updatedAtDesc"
)
