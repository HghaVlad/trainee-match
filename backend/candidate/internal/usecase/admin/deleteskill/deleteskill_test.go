package deleteskill

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

var errDB = errors.New("db error")

func TestExecute(t *testing.T) {
	ctx := context.Background()
	skillID := uuid.New()

	tests := []struct {
		name          string
		mockSetup     func(repo *MockSkillRepo)
		expectedError error
	}{
		{
			name: "delete ok",
			mockSetup: func(repo *MockSkillRepo) {
				repo.On("Delete", ctx, skillID).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "not found",
			mockSetup: func(repo *MockSkillRepo) {
				repo.On("Delete", ctx, skillID).Return(domain.ErrSkillNotFound).Once()
			},
			expectedError: domain.ErrSkillNotFound,
		},
		{
			name: "repo error",
			mockSetup: func(repo *MockSkillRepo) {
				repo.On("Delete", ctx, skillID).Return(errDB).Once()
			},
			expectedError: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockSkillRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			uc := NewUseCase(repo)
			err := uc.Execute(ctx, Request{SkillID: skillID})
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
