package handlers

import (
	"context"
	"log/slog"

	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/apply"
	appviews "github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/identity"
)

type Deps struct {
	Apply  applyUC
	Logger *slog.Logger
}

type applyUC interface {
	Execute(context.Context, apply.Request, identity.Identity) (*appviews.Details, error)
}
