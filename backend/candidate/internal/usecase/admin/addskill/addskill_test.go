package addskill

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
	validName := "Golang"
	validID := uuid.New()

	tests := []struct {
		name          string
		request       Request
		mockSetup     func(repo *MockSkillRepo)
		expectedID    uuid.UUID
		expectedError error
	}{
		{
			name:    "valid create",
			request: Request{Name: validName},
			mockSetup: func(repo *MockSkillRepo) {
				repo.On("Create", ctx, domain.Skill{Name: validName}).Return(validID, nil).Once()
			},
			expectedID:    validID,
			expectedError: nil,
		},
		{
			name:          "invalid name",
			request:       Request{Name: ""},
			mockSetup:     nil,
			expectedID:    uuid.Nil,
			expectedError: domain.ErrInvalidSkillName,
		},
		{
			name:    "repo error",
			request: Request{Name: validName},
			mockSetup: func(repo *MockSkillRepo) {
				repo.On("Create", ctx, domain.Skill{Name: validName}).Return(uuid.Nil, errDB).Once()
			},
			expectedID:    uuid.Nil,
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
			resp, err := uc.Execute(ctx, tt.request)
			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
				require.Equal(t, uuid.Nil, resp.ID)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedID, resp.ID)
				require.Equal(t, validName, resp.Name)
			}

			repo.AssertExpectations(t)
		})
	}
}
