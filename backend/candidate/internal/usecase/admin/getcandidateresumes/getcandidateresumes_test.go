package getcandidateresumes

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
	candidateID := uuid.New()
	page := 1
	size := 2

	resumes := []domain.Resume{
		{
			ID:               uuid.New(),
			CandidateId:      candidateID,
			Name:             "Resume 1",
			Status:           0,
			ModerationStatus: domain.ModerationStatusOK,
		},
		{
			ID:               uuid.New(),
			CandidateId:      candidateID,
			Name:             "Resume 2",
			Status:           1,
			ModerationStatus: domain.ModerationStatusHidden,
		},
	}

	tests := []struct {
		name          string
		mockSetup     func(repo *MockResumeRepo)
		expectedCount int
		expectedError error
	}{
		{
			name: "valid list",
			mockSetup: func(repo *MockResumeRepo) {
				repo.On("GetByCandidateId", ctx, candidateID, page, size).Return(resumes, nil).Once()
			},
			expectedCount: len(resumes),
			expectedError: nil,
		},
		{
			name: "invalid status",
			mockSetup: func(repo *MockResumeRepo) {
				bad := []domain.Resume{{ID: uuid.New(), CandidateId: candidateID, Name: "Bad", Status: 99}}
				repo.On("GetByCandidateId", ctx, candidateID, page, size).Return(bad, nil).Once()
			},
			expectedCount: 0,
			expectedError: domain.ErrInvalidResumeStatus,
		},
		{
			name: "repo error",
			mockSetup: func(repo *MockResumeRepo) {
				repo.On("GetByCandidateId", ctx, candidateID, page, size).Return(nil, errDB).Once()
			},
			expectedCount: 0,
			expectedError: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockResumeRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			uc := NewUseCase(repo)
			resp, err := uc.Execute(ctx, Request{CandidateID: candidateID, Page: page, Size: size})
			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
				require.Len(t, resp, tt.expectedCount)
				require.Equal(t, resumes[0].ID, resp[0].ID)
				require.Equal(t, resumes[0].CandidateId, resp[0].CandidateID)
				require.Equal(t, resumes[0].Name, resp[0].Name)
				require.Equal(t, "draft", resp[0].Status)
				require.Equal(t, string(resumes[0].ModerationStatus), resp[0].ModerationStatus)
			}

			repo.AssertExpectations(t)
		})
	}
}
