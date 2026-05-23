package remove_resume

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

var (
	ErrResumeDb    = errors.New("resume db error")
	ErrCandidateDb = errors.New("candidate db error")
)

func TestExecute(t *testing.T) {
	ctx := context.Background()
	resumeID := uuid.New()
	userID := uuid.New()
	candidateID := uuid.New()

	tests := []struct {
		name          string
		req           Request
		mockSetup     func(candidateRepo *MockCandidateRepo, resumeRepo *MockResumeRepo)
		expectedError error
	}{
		{
			name: "success",
			req:  Request{ResumeId: resumeID, UserId: userID},
			mockSetup: func(candidateRepo *MockCandidateRepo, resumeRepo *MockResumeRepo) {
				candidateRepo.On("GetByUserID", ctx, userID).
					Return(domain.Candidate{ID: candidateID, UserId: userID}, nil).
					Once()
				resumeRepo.On("Remove", ctx, resumeID, candidateID).
					Return(nil).
					Once()
			},
			expectedError: nil,
		},
		{
			name: "candidate not found",
			req:  Request{ResumeId: resumeID, UserId: userID},
			mockSetup: func(candidateRepo *MockCandidateRepo, resumeRepo *MockResumeRepo) {
				candidateRepo.On("GetByUserID", ctx, userID).
					Return(domain.Candidate{}, domain.ErrCandidateNotFound).
					Once()
			},
			expectedError: domain.ErrCandidateNotFound,
		},
		{
			name: "candidate repo error",
			req:  Request{ResumeId: resumeID, UserId: userID},
			mockSetup: func(candidateRepo *MockCandidateRepo, resumeRepo *MockResumeRepo) {
				candidateRepo.On("GetByUserID", ctx, userID).
					Return(domain.Candidate{}, ErrCandidateDb).
					Once()
			},
			expectedError: ErrCandidateDb,
		},
		{
			name: "resume not found",
			req:  Request{ResumeId: resumeID, UserId: userID},
			mockSetup: func(candidateRepo *MockCandidateRepo, resumeRepo *MockResumeRepo) {
				candidateRepo.On("GetByUserID", ctx, userID).
					Return(domain.Candidate{ID: candidateID, UserId: userID}, nil).
					Once()
				resumeRepo.On("Remove", ctx, resumeID, candidateID).
					Return(domain.ErrResumeNotFound).
					Once()
			},
			expectedError: domain.ErrResumeNotFound,
		},
		{
			name: "resume repo error",
			req:  Request{ResumeId: resumeID, UserId: userID},
			mockSetup: func(candidateRepo *MockCandidateRepo, resumeRepo *MockResumeRepo) {
				candidateRepo.On("GetByUserID", ctx, userID).
					Return(domain.Candidate{ID: candidateID, UserId: userID}, nil).
					Once()
				resumeRepo.On("Remove", ctx, resumeID, candidateID).
					Return(ErrResumeDb).
					Once()
			},
			expectedError: ErrResumeDb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			candidateRepo := &MockCandidateRepo{}
			resumeRepo := &MockResumeRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(candidateRepo, resumeRepo)
			}

			writer := &MockEventWriter{}
			writer.On("WriteResumeDeleted", context.Background(), mock.AnythingOfType("events.ResumeDeleted")).
				Return(nil)

			trManager := &MockTrManager{}
			trManager.On("Do", mock.Anything, mock.Anything).
				Return(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})

			uc := New(resumeRepo, candidateRepo, writer, trManager)
			err := uc.RemoveResume(ctx, tt.req)

			if tt.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
			}

			candidateRepo.AssertExpectations(t)
			resumeRepo.AssertExpectations(t)
		})
	}
}
