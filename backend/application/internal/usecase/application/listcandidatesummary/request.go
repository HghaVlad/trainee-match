package listcandidatesummary

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/cursors"
)

type Request struct {
	Statuses  []application.Status
	CompanyID *uuid.UUID
	Cursor    string
	Limit     int
	Order     cursors.SummaryOrder
}
