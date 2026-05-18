package me

import (
	"context"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/views"
)

type repo interface {
	GetFullView(ctx context.Context, userID, companyID uuid.UUID) (*views.MemberFullView, error)
}
