package getcandidates

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
	page := 1
	size := 2

	candidates := []domain.Candidate{
		{
			ID:       uuid.New(),
			FullName: "John Doe",
			UserId:   uuid.New(),
			Phone:    "+1234567890",
			Telegram: "@john",
			City:     "Kyiv",
			Birthday: time.Date(1995, time.May, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:       uuid.New(),
			FullName: "Jane Smith",
			UserId:   uuid.New(),
			Phone:    "+1098765432",
			Telegram: "@jane",
			City:     "Lviv",
			Birthday: time.Date(1998, time.January, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	tests := []struct {
		name          string
		mockSetup     func(repo *MockCandidateRepo)
		expectedCount int
		expectedError error
	}{
		{
			name: "valid list",
			mockSetup: func(repo *MockCandidateRepo) {
				repo.On("GetCandidates", ctx, page, size).Return(candidates, nil).Once()
			},
			expectedCount: len(candidates),
			expectedError: nil,
		},
		{
			name: "repo error",
			mockSetup: func(repo *MockCandidateRepo) {
				repo.On("GetCandidates", ctx, page, size).Return(nil, errDB).Once()
			},
			expectedCount: 0,
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
			resp, err := uc.Execute(ctx, Request{Page: page, Size: size})
			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.Len(t, resp, tt.expectedCount)
				require.Equal(t, candidates[0].ID, resp[0].ID)
				require.Equal(t, candidates[0].FullName, resp[0].FullName)
				require.Equal(t, candidates[0].UserId, resp[0].UserID)
				require.Equal(t, candidates[0].Phone, resp[0].Phone)
				require.Equal(t, candidates[0].Telegram, resp[0].Telegram)
				require.Equal(t, candidates[0].City, resp[0].City)
				require.WithinDuration(t, candidates[0].Birthday, resp[0].Birthday, time.Second)
			}

			repo.AssertExpectations(t)
		})
	}
}
