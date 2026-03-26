package get_resume

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/get_resume/mocks"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
)

var (
	ErrDb = errors.New("db error")
)

func TestGetById(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	candidateId := uuid.New()
	now := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	domainResume := domain.Resume{
		ID:          id,
		CandidateId: candidateId,
		Name:        "Resume 1",
		Status:      domain.Draft,
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

	validResp := &Response{
		ID:          id,
		CandidateID: candidateId,
		Name:        domainResume.Name,
		Status:      domainResume.Status,
		Data:        convertDomainDataToResponseData(domainResume.Data),
	}

	tests := []struct {
		name          string
		mockSetup     func(resumeRepo *mocks.ResumeRepo)
		reqUserId     uuid.UUID
		expectedError error
	}{
		{
			name: "valid get",
			mockSetup: func(resumeRepo *mocks.ResumeRepo) {
				resumeRepo.On("GetById", ctx, id).Return(domainResume, nil).Once()
			},
			reqUserId:     uuid.New(),
			expectedError: nil,
		},
		{
			name: "not found",
			mockSetup: func(resumeRepo *mocks.ResumeRepo) {
				resumeRepo.On("GetById", ctx, id).Return(domain.Resume{}, domain.ErrResumeNotFound).Once()
			},
			reqUserId:     uuid.New(),
			expectedError: domain.ErrResumeNotFound,
		},
		{
			name: "repo error",
			mockSetup: func(resumeRepo *mocks.ResumeRepo) {
				resumeRepo.On("GetById", ctx, id).Return(domain.Resume{}, ErrDb).Once()
			},
			reqUserId:     uuid.New(),
			expectedError: ErrDb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resumeRepo := &mocks.ResumeRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(resumeRepo)
			}

			uc := New(resumeRepo, nil)
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
		{ID: uuid.New(), CandidateId: candidateId, Name: "r1", Status: domain.Draft},
		{ID: uuid.New(), CandidateId: candidateId, Name: "r2", Status: domain.Published},
	}

	tests := []struct {
		name          string
		mockSetup     func(resumeRepo *mocks.ResumeRepo, candidateRepo *mocks.CandidateRepo)
		expectedError error
	}{
		{
			name: "valid get",
			mockSetup: func(resumeRepo *mocks.ResumeRepo, candidateRepo *mocks.CandidateRepo) {
				candidateRepo.On("GetByUserID", ctx, userId).Return(domain.Candidate{ID: candidateId, UserId: userId}, nil).Once()
				resumeRepo.On("GetByCandidateId", ctx, candidateId).Return(domainResumes, nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "candidate not found",
			mockSetup: func(resumeRepo *mocks.ResumeRepo, candidateRepo *mocks.CandidateRepo) {
				candidateRepo.On("GetByUserID", ctx, userId).Return(domain.Candidate{}, domain.ErrCandidateNotFound).Once()
			},
			expectedError: domain.ErrCandidateNotFound,
		},
		{
			name: "repo error",
			mockSetup: func(resumeRepo *mocks.ResumeRepo, candidateRepo *mocks.CandidateRepo) {
				candidateRepo.On("GetByUserID", ctx, userId).Return(domain.Candidate{ID: candidateId, UserId: userId}, nil).Once()
				resumeRepo.On("GetByCandidateId", ctx, candidateId).Return(nil, ErrDb).Once()
			},
			expectedError: ErrDb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resumeRepo := &mocks.ResumeRepo{}
			candidateRepo := &mocks.CandidateRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(resumeRepo, candidateRepo)
			}

			uc := New(resumeRepo, candidateRepo)
			resp, err := uc.GetByCandidateId(ctx, userId)

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
