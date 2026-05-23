package application_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
)

func TestApplication_ChangeStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		initialStatus application.Status
		nextStatus    application.Status
		actor         application.Actor
		wantErr       error
	}{
		{
			name:          "submitted -> seen by hr",
			initialStatus: application.StatusSubmitted,
			nextStatus:    application.StatusSeen,
			actor:         application.ActorHR,
			wantErr:       nil,
		},
		{
			name:          "submitted -> interview by hr",
			initialStatus: application.StatusSubmitted,
			nextStatus:    application.StatusInterview,
			actor:         application.ActorHR,
			wantErr:       nil,
		},
		{
			name:          "submitted -> withdrawn by candidate",
			initialStatus: application.StatusSubmitted,
			nextStatus:    application.StatusWithdrawn,
			actor:         application.ActorCandidate,
			wantErr:       nil,
		},
		{
			name:          "submitted -> offer invalid",
			initialStatus: application.StatusSubmitted,
			nextStatus:    application.StatusOffer,
			actor:         application.ActorHR,
			wantErr:       application.ErrInvalidStatusTransition,
		},
		{
			name:          "submitted -> seen by candidate invalid",
			initialStatus: application.StatusSubmitted,
			nextStatus:    application.StatusSeen,
			actor:         application.ActorCandidate,
			wantErr:       application.ErrInvalidStatusTransition,
		},
		{
			name:          "same status",
			initialStatus: application.StatusSubmitted,
			nextStatus:    application.StatusSubmitted,
			actor:         application.ActorCandidate,
			wantErr:       application.ErrStatusAlreadySet,
		},
		{
			name:          "rejected cannot transition",
			initialStatus: application.StatusRejected,
			nextStatus:    application.StatusSeen,
			actor:         application.ActorHR,
			wantErr:       application.ErrInvalidStatusTransition,
		},
		{
			name:          "withdrawn cannot transition",
			initialStatus: application.StatusWithdrawn,
			nextStatus:    application.StatusInterview,
			actor:         application.ActorHR,
			wantErr:       application.ErrInvalidStatusTransition,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			now := time.Now().UTC()

			app := &application.Application{
				ID:        uuid.New(),
				Status:    tt.initialStatus,
				CreatedAt: now,
				UpdatedAt: now,
			}

			nextTime := now.Add(time.Hour)

			err := app.ChangeStatus(
				tt.nextStatus,
				tt.actor,
				nextTime,
			)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.nextStatus, app.Status)
			require.Equal(t, nextTime, app.UpdatedAt)
		})
	}
}
