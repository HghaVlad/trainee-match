package me

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/views"
	"github.com/google/uuid"
)

type repo interface {
	GetFullView(ctx context.Context, userID, companyID uuid.UUID) (*views.MemberFullView, error)
}
