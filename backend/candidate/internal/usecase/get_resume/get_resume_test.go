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

	validResp := &GetByIdResponse{
		ID:          id,
		CandidateId: candidateId,
		Name:        domainResume.Name,
		Status:      domainResume.Status,
		Data:        convertDomainDataToResponseData(domainResume.Data),
	}

	tests := []struct {
		name          string
		mockSetup     func(repo *mocks.ResumeRepo)
		expectedError error
	}{
		{
			name: "valid get",
			mockSetup: func(repo *mocks.ResumeRepo) {
				repo.On("GetById", ctx, id).Return(domainResume, nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "not found",
			mockSetup: func(repo *mocks.ResumeRepo) {
				repo.On("GetById", ctx, id).Return(domain.Resume{}, domain.ErrResumeNotFound).Once()
			},
			expectedError: domain.ErrResumeNotFound,
		},
		{
			name: "repo error",
			mockSetup: func(repo *mocks.ResumeRepo) {
				repo.On("GetById", ctx, id).Return(domain.Resume{}, ErrDb).Once()
			},
			expectedError: ErrDb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.ResumeRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			uc := New(repo)
			resp, err := uc.GetById(ctx, GetByIdRequest{ID: id})
			if tt.expectedError != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.expectedError), "expected error: %v, got: %v", tt.expectedError, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, validResp.ID, resp.ID)
				require.Equal(t, validResp.CandidateId, resp.CandidateId)
				require.Equal(t, validResp.Name, resp.Name)
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestGetByCandidateId(t *testing.T) {
	ctx := context.Background()
	candidateId := uuid.New()

	domainResumes := []domain.Resume{
		{ID: uuid.New(), CandidateId: candidateId, Name: "r1", Status: domain.Draft},
		{ID: uuid.New(), CandidateId: candidateId, Name: "r2", Status: domain.Published},
	}

	tests := []struct {
		name          string
		mockSetup     func(repo *mocks.ResumeRepo)
		expectedError error
	}{
		{
			name: "valid get",
			mockSetup: func(repo *mocks.ResumeRepo) {
				repo.On("GetByCandidateId", ctx, candidateId).Return(domainResumes, nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "repo error",
			mockSetup: func(repo *mocks.ResumeRepo) {
				repo.On("GetByCandidateId", ctx, candidateId).Return(nil, ErrDb).Once()
			},
			expectedError: ErrDb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.ResumeRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			uc := New(repo)
			resp, err := uc.GetByCandidateId(ctx, GetByCandidateIdRequest{CandidateId: candidateId})

			if tt.expectedError != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.expectedError), "expected error: %v, got: %v", tt.expectedError, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, len(domainResumes), len(resp))
			}

			repo.AssertExpectations(t)
		})
	}
}
