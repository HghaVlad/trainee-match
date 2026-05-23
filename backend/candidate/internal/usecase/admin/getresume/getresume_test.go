package getresume

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
	resumeID := uuid.New()
	candidateID := uuid.New()

	resume := domain.Resume{
		ID:               resumeID,
		CandidateId:      candidateID,
		Name:             "Resume",
		Status:           1,
		ModerationStatus: domain.ModerationStatusOK,
		Data: domain.ResumeData{
			FirstName:   "John",
			LastName:    "Doe",
			DateOfBirth: time.Date(1996, time.June, 10, 0, 0, 0, 0, time.UTC),
			Email:       "john@example.com",
			Phone:       "+1234567890",
			City:        "Kyiv",
			Citizenship: "UA",
		},
	}

	tests := []struct {
		name          string
		mockSetup     func(repo *MockResumeRepo)
		expectedError error
	}{
		{
			name: "valid get",
			mockSetup: func(repo *MockResumeRepo) {
				repo.On("GetById", ctx, resumeID).Return(resume, nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "invalid status",
			mockSetup: func(repo *MockResumeRepo) {
				bad := resume
				bad.Status = 99
				repo.On("GetById", ctx, resumeID).Return(bad, nil).Once()
			},
			expectedError: domain.ErrInvalidResumeStatus,
		},
		{
			name: "repo error",
			mockSetup: func(repo *MockResumeRepo) {
				repo.On("GetById", ctx, resumeID).Return(domain.Resume{}, errDB).Once()
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

			uc := NewUseCase(repo)
			resp, err := uc.Execute(ctx, Request{ResumeID: resumeID})
			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, resume.ID, resp.ID)
				require.Equal(t, resume.CandidateId, resp.CandidateID)
				require.Equal(t, resume.Name, resp.Name)
				require.Equal(t, "published", resp.Status)
				require.Equal(t, string(resume.ModerationStatus), resp.ModerationStatus)
				require.Equal(t, resume.Data.FirstName, resp.Data.FirstName)
				require.Equal(t, resume.Data.LastName, resp.Data.LastName)
				require.WithinDuration(t, resume.Data.DateOfBirth, resp.Data.DateOfBirth, time.Second)
			}

			repo.AssertExpectations(t)
		})
	}
}
