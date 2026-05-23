package getcandidate

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

var errDB = errors.New("db error")

func TestExecute(t *testing.T) {
	ctx := context.Background()
	candidateID := uuid.New()
	birthday := time.Date(1994, time.December, 10, 0, 0, 0, 0, time.UTC)

	dCandidate := domain.Candidate{
		ID:       candidateID,
		FullName: "John Doe",
		UserId:   uuid.New(),
		Phone:    "+1234567890",
		Telegram: "@john",
		City:     "Kyiv",
		Birthday: birthday,
	}

	tests := []struct {
		name          string
		mockSetup     func(repo *MockCandidateRepo)
		expectedError error
	}{
		{
			name: "valid get",
			mockSetup: func(repo *MockCandidateRepo) {
				repo.On("GetByID", ctx, candidateID).Return(dCandidate, nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "not found",
			mockSetup: func(repo *MockCandidateRepo) {
				repo.On("GetByID", ctx, candidateID).Return(domain.Candidate{}, domain.ErrCandidateNotFound).Once()
			},
			expectedError: domain.ErrCandidateNotFound,
		},
		{
			name: "repo error",
			mockSetup: func(repo *MockCandidateRepo) {
				repo.On("GetByID", ctx, candidateID).Return(domain.Candidate{}, errDB).Once()
			},
			expectedError: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockCandidateRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			uc := NewUseCase(repo)
			resp, err := uc.Execute(ctx, Request{CandidateID: candidateID})
			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, dCandidate.ID, resp.ID)
				require.Equal(t, dCandidate.FullName, resp.FullName)
				require.Equal(t, dCandidate.UserId, resp.UserID)
				require.Equal(t, dCandidate.Phone, resp.Phone)
				require.Equal(t, dCandidate.Telegram, resp.Telegram)
				require.Equal(t, dCandidate.City, resp.City)
				require.WithinDuration(t, dCandidate.Birthday, resp.Birthday, time.Second)
			}

			repo.AssertExpectations(t)
		})
	}
}
