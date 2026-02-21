package update_resume

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/update_resume/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
)

// TODO: add more tests

var (
	ErrUpdateDb = errors.New("update db error")
)

func TestExecute(t *testing.T) {
	ctx := context.Background()
	resumeID := uuid.New()
	candidateID := uuid.New()

	existing := domain.Resume{
		ID:          resumeID,
		CandidateId: candidateID,
		Name:        "Old",
		Status:      domain.Draft,
		Data: domain.ResumeData{
			LastName: "Doe", FirstName: "John", DateOfBirth: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
			Email: "john@example.com", Phone: "+1234567890", City: "City", Citizenship: "Country",
			Education:       []domain.Education{{Level: "BSc", University: "Uni", StartYear: 2008, EndYear: 2012}},
			WorkExperiences: []domain.WorkExperience{{Position: "Dev", Company: "Co", Period: "2012-2018"}},
			SkillsList:      []uuid.UUID{uuid.New()},
		},
	}

	tests := []struct {
		name          string
		req           Request
		mockSetup     func(repo *mocks.ResumeRepo, skillRepo *mocks.SkillRepo)
		expectedError error
	}{
		{
			name: "valid update",
			req:  Request{ID: resumeID, Name: stringPtr("New Name")},
			mockSetup: func(repo *mocks.ResumeRepo, skillRepo *mocks.SkillRepo) {
				repo.On("GetById", ctx, resumeID).Return(existing, nil).Once()
				repo.On("Update", ctx, mock.AnythingOfType("*domain.Resume")).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "not found",
			req:  Request{ID: uuid.New()},
			mockSetup: func(repo *mocks.ResumeRepo, skillRepo *mocks.SkillRepo) {
				repo.On("GetById", ctx, mock.Anything).Return(domain.Resume{}, domain.ErrResumeNotFound).Once()
			},
			expectedError: domain.ErrResumeNotFound,
		},
		{
			name: "skills not exist",
			req:  Request{ID: resumeID, Data: &ResumeData{SkillsList: &[]uuid.UUID{uuid.New()}}},
			mockSetup: func(repo *mocks.ResumeRepo, skillRepo *mocks.SkillRepo) {
				repo.On("GetById", ctx, resumeID).Return(existing, nil).Once()
				skillRepo.On("AreSkillsExist", ctx, mock.Anything).Return(false, nil).Once()
			},
			expectedError: domain.ErrSkillNotFound,
		},
		{
			name: "validation fail after update",
			req:  Request{ID: resumeID, Data: &ResumeData{Phone: stringPtr("bad phone")}},
			mockSetup: func(repo *mocks.ResumeRepo, skillRepo *mocks.SkillRepo) {
				repo.On("GetById", ctx, resumeID).Return(existing, nil).Once()
			},
			expectedError: domain.ErrInvalidPhoneFormat,
		},
		{
			name: "repo update error",
			req:  Request{ID: resumeID, Name: stringPtr("New")},
			mockSetup: func(repo *mocks.ResumeRepo, skillRepo *mocks.SkillRepo) {
				repo.On("GetById", ctx, resumeID).Return(existing, nil).Once()
				repo.On("Update", ctx, mock.AnythingOfType("*domain.Resume")).Return(ErrUpdateDb).Once()
			},
			expectedError: ErrUpdateDb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.ResumeRepo{}
			skillRepo := &mocks.SkillRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(repo, skillRepo)
			}

			uc := New(repo, skillRepo)
			res, err := uc.Execute(ctx, tt.req)
			if tt.expectedError != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.expectedError), "expected %v got %v", tt.expectedError, err)
			} else {
				require.NoError(t, err)
				require.True(t, res.Success)
			}

			repo.AssertExpectations(t)
			skillRepo.AssertExpectations(t)
		})
	}
}

func stringPtr(s string) *string { return &s }
