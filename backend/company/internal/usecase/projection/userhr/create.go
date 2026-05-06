package userhr

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
)

type CreatedEvent struct {
	EventID    uuid.UUID           `avro:"event_id"`
	UserID     uuid.UUID           `avro:"user_id"`
	Username   string              `avro:"username"`
	Role       identity.GlobalRole `avro:"role"`
	Email      string              `avro:"email"`
	OccurredAt time.Time           `avro:"occurred_at"`
}

type createRepo interface {
	CreateIdempotent(ctx context.Context, hrProj Projection) error
}

type CreateUsecase struct {
	createRepo createRepo
}

func NewCreatedUsecase(createRepo createRepo) *CreateUsecase {
	return &CreateUsecase{createRepo: createRepo}
}

func (uc *CreateUsecase) Execute(ctx context.Context, ev CreatedEvent) error {
	if ev.Role != identity.RoleHR {
		return nil
	}

	proj, err := NewHrProjection(ev.UserID, ev.Username, ev.Email, ev.OccurredAt)
	if err != nil {
		return err
	}

	err = uc.createRepo.CreateIdempotent(ctx, *proj)
	if err != nil {
		return err
	}

	return nil
}
