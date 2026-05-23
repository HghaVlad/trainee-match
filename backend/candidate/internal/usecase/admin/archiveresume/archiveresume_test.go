package archiveresume

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

var errDB = errors.New("db error")

func TestExecute(t *testing.T) {
	ctx := context.Background()
	resumeID := uuid.New()

	tests := []struct {
		name          string
		mockSetup     func(repo *MockResumeRepo)
		expectedError error
	}{
		{
			name: "already hidden",
			mockSetup: func(repo *MockResumeRepo) {
				repo.On("GetById", ctx, resumeID).
					Return(domain.Resume{ID: resumeID, ModerationStatus: domain.ModerationStatusHidden}, nil).
					Once()
			},
			expectedError: nil,
		},
		{
			name: "set hidden",
			mockSetup: func(repo *MockResumeRepo) {
				repo.On("GetById", ctx, resumeID).
					Return(domain.Resume{ID: resumeID, ModerationStatus: domain.ModerationStatusOK}, nil).
					Once()
				repo.On("SetModerationStatus", ctx, resumeID, domain.ModerationStatusHidden).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "get error",
			mockSetup: func(repo *MockResumeRepo) {
				repo.On("GetById", ctx, resumeID).Return(domain.Resume{}, errDB).Once()
			},
			expectedError: errDB,
		},
		{
			name: "update error",
			mockSetup: func(repo *MockResumeRepo) {
				repo.On("GetById", ctx, resumeID).
					Return(domain.Resume{ID: resumeID, ModerationStatus: domain.ModerationStatusOK}, nil).
					Once()
				repo.On("SetModerationStatus", ctx, resumeID, domain.ModerationStatusHidden).Return(errDB).Once()
			},
			expectedError: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockResumeRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			writer := &MockEventWriter{}
			writer.On("WriteResumeArchived", context.Background(), mock.AnythingOfType("events.ResumeArchived")).
				Return(nil)

			trManager := &MockTrManager{}
			trManager.On("Do", mock.Anything, mock.Anything).
				Return(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})

			uc := NewUseCase(repo, writer, trManager)
			err := uc.Execute(ctx, Request{ResumeID: resumeID})
			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}
