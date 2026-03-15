package get_candidate

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/get_candidate/mocks"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
)

var (
	ErrDB = errors.New("db error")
)

func TestExecute(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	userId := uuid.New()
	birthday := time.Date(1995, time.May, 10, 0, 0, 0, 0, time.UTC)

	dCandidate := domain.Candidate{
		ID:       id,
		UserId:   userId,
		Phone:    "+1234567890",
		Telegram: "@user",
		City:     "City",
		Birthday: birthday,
	}

	tests := []struct {
		name          string
		mockSetup     func(repo *mocks.CandidateRepo)
		expectedError error
	}{
		{
			name:          "valid get",
			mockSetup:     func(repo *mocks.CandidateRepo) { repo.On("GetByID", ctx, id).Return(dCandidate, nil).Once() },
			expectedError: nil,
		},
		{
			name: "not found",
			mockSetup: func(repo *mocks.CandidateRepo) {
				repo.On("GetByID", ctx, id).Return(domain.Candidate{}, domain.ErrCandidateNotFound).Once()
			},
			expectedError: domain.ErrCandidateNotFound,
		},
		{
			name: "repo error",
			mockSetup: func(repo *mocks.CandidateRepo) {
				repo.On("GetByID", ctx, id).Return(domain.Candidate{}, ErrDB).Once()
			},
			expectedError: ErrDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.CandidateRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			uc := New(repo)
			resp, err := uc.Execute(ctx, id)
			if tt.expectedError != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.expectedError))
			} else {
				require.NoError(t, err)
				require.Equal(t, dCandidate.ID, resp.ID)
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
