package get_resume

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

var (
	ErrDb = errors.New("db error")
)

func TestGetById(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	userID := uuid.New()
	candidateId := uuid.New()
	now := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	domainResume := domain.Resume{
		ID:          id,
		CandidateId: candidateId,
		Name:        "Resume 1",
		Status:      0,
		Data: domain.ResumeData{
			LastName:        "Doe",
			FirstName:       "John",
			DateOfBirth:     now,
			Email:           "john@example.com",
			Phone:           "+1234567890",
			City:            "City",
			Citizenship:     "Country",
			Education:       []domain.Education{{Level: "BSc", University: "Uni", StartYear: 2008, EndYear: 2012}},
			WorkExperiences: []domain.WorkExperience{{Position: "Dev", Company: "Co", Period: "2012-2018"}},
			SkillsList:      []uuid.UUID{uuid.New()},
		},
	}
	domainCandidate := domain.Candidate{
		ID:     candidateId,
		UserId: userID,
	}

	validResp := &Response{
		ID:          id,
		CandidateID: candidateId,
		Name:        domainResume.Name,
		Status:      "draft",
		Data:        convertDomainDataToResponseData(domainResume.Data),
	}

	tests := []struct {
		name          string
		mockSetup     func(resumeRepo *MockResumeRepo, candidateRepo *MockCandidateRepo)
		reqUserId     uuid.UUID
		expectedError error
	}{
		{
			name: "valid get",
			mockSetup: func(resumeRepo *MockResumeRepo, candidateRepo *MockCandidateRepo) {
				resumeRepo.On("GetById", ctx, id).Return(domainResume, nil).Once()
				candidateRepo.On("GetByUserID", ctx, userID).Return(domainCandidate, nil).Once()
			},
			reqUserId:     userID,
			expectedError: nil,
		},
		{
			name: "not found",
			mockSetup: func(resumeRepo *MockResumeRepo, candidateRepo *MockCandidateRepo) {
				candidateRepo.On("GetByUserID", ctx, userID).Return(domainCandidate, nil).Once()
				resumeRepo.On("GetById", ctx, id).Return(domain.Resume{}, domain.ErrResumeNotFound).Once()
			},
			reqUserId:     userID,
			expectedError: domain.ErrResumeNotFound,
		},
		{
			name: "another user",
			mockSetup: func(resumeRepo *MockResumeRepo, candidateRepo *MockCandidateRepo) {
				resumeRepo.On("GetById", ctx, id).Return(domain.Resume{ID: id, CandidateId: uuid.New()}, nil).Once()
				candidateRepo.On("GetByUserID", ctx, userID).Return(domainCandidate, nil).Once()
			},
			reqUserId:     userID,
			expectedError: domain.ErrResumeNotFound,
		},
		{
			name: "resume repo error",
			mockSetup: func(resumeRepo *MockResumeRepo, candidateRepo *MockCandidateRepo) {
				candidateRepo.On("GetByUserID", ctx, userID).Return(domainCandidate, nil).Once()
				resumeRepo.On("GetById", ctx, id).Return(domain.Resume{}, ErrDb).Once()
			},
			reqUserId:     userID,
			expectedError: ErrDb,
		},
		{
			name: "candidate repo error",
			mockSetup: func(resumeRepo *MockResumeRepo, candidateRepo *MockCandidateRepo) {
				resumeRepo.On("GetById", ctx, id).Return(domainResume, nil).Maybe()
				candidateRepo.On("GetByUserID", ctx, userID).Return(domain.Candidate{}, ErrDb).Once()
			},
			reqUserId:     userID,
			expectedError: ErrDb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resumeRepo := &MockResumeRepo{}
			candidateRepo := &MockCandidateRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(resumeRepo, candidateRepo)
			}

			uc := New(resumeRepo, candidateRepo)
			resp, err := uc.GetById(ctx, id, tt.reqUserId)
			if tt.expectedError != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.expectedError), "expected error: %v, got: %v", tt.expectedError, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, validResp.ID, resp.ID)
				require.Equal(t, validResp.CandidateID, resp.CandidateID)
				require.Equal(t, validResp.Name, resp.Name)
			}

			resumeRepo.AssertExpectations(t)
		})
	}
}

func TestGetByCandidateId(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	candidateId := uuid.New()
	domainResumes := []domain.Resume{
		{ID: uuid.New(), CandidateId: candidateId, Name: "r1", Status: 0},
		{ID: uuid.New(), CandidateId: candidateId, Name: "r2", Status: 1},
	}

	tests := []struct {
		name          string
		mockSetup     func(resumeRepo *MockResumeRepo, candidateRepo *MockCandidateRepo)
		expectedError error
	}{
		{
			name: "valid get",
			mockSetup: func(resumeRepo *MockResumeRepo, candidateRepo *MockCandidateRepo) {
				candidateRepo.On("GetByUserID", ctx, userId).
					Return(domain.Candidate{ID: candidateId, UserId: userId}, nil).
					Once()
				resumeRepo.On("GetByCandidateId", ctx, candidateId, 1, 10).Return(domainResumes, nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "candidate not found",
			mockSetup: func(resumeRepo *MockResumeRepo, candidateRepo *MockCandidateRepo) {
				candidateRepo.On("GetByUserID", ctx, userId).
					Return(domain.Candidate{}, domain.ErrCandidateNotFound).
					Once()
			},
			expectedError: domain.ErrCandidateNotFound,
		},
		{
			name: "repo error",
			mockSetup: func(resumeRepo *MockResumeRepo, candidateRepo *MockCandidateRepo) {
				candidateRepo.On("GetByUserID", ctx, userId).
					Return(domain.Candidate{ID: candidateId, UserId: userId}, nil).
					Once()
				resumeRepo.On("GetByCandidateId", ctx, candidateId, 1, 10).Return(nil, ErrDb).Once()
			},
			expectedError: ErrDb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resumeRepo := &MockResumeRepo{}
			candidateRepo := &MockCandidateRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(resumeRepo, candidateRepo)
			}

			uc := New(resumeRepo, candidateRepo)
			resp, err := uc.GetByCandidateId(ctx, userId, 1, 10)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.expectedError), "expected error: %v, got: %v", tt.expectedError, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, len(domainResumes), len(resp))
			}

			resumeRepo.AssertExpectations(t)
			candidateRepo.AssertExpectations(t)
		})
	}
}
