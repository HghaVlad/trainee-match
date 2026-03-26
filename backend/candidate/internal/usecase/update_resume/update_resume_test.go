package update_resume

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/update_resume/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TODO: add more tests

var (
	ErrUpdateDb    = errors.New("update db error")
	ErrCandidateDb = errors.New("candidate db error")
)

func TestExecute(t *testing.T) {
	ctx := context.Background()
	resumeID := uuid.New()
	candidateID := uuid.New()
	userId := uuid.New()
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
	candidate := domain.Candidate{
		ID:     candidateID,
		UserId: userId,
	}

	tests := []struct {
		name          string
		req           Request
		mockSetup     func(repo *mocks.ResumeRepo, candidateRepo *mocks.CandidateRepo, skillRepo *mocks.SkillRepo)
		expectedError error
	}{
		{
			name: "valid update",
			req:  Request{ID: resumeID, UserId: userId, Name: stringPtr("New Name")},
			mockSetup: func(repo *mocks.ResumeRepo, candidateRepo *mocks.CandidateRepo, skillRepo *mocks.SkillRepo) {
				candidateRepo.On("GetByUserID", ctx, userId).Return(candidate, nil).Once()

				repo.On("GetById", ctx, resumeID).Return(existing, nil).Once()
				repo.On("Update", ctx, mock.AnythingOfType("*domain.Resume")).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "not found",
			req:  Request{ID: uuid.New(), UserId: userId},
			mockSetup: func(repo *mocks.ResumeRepo, candidateRepo *mocks.CandidateRepo, skillRepo *mocks.SkillRepo) {
				candidateRepo.On("GetByUserID", ctx, userId).Return(candidate, nil).Once()

				repo.On("GetById", ctx, mock.Anything).Return(domain.Resume{}, domain.ErrResumeNotFound).Once()
			},
			expectedError: domain.ErrResumeNotFound,
		},
		{
			name: "forbidden",
			req:  Request{ID: resumeID, UserId: userId},
			mockSetup: func(repo *mocks.ResumeRepo, candidateRepo *mocks.CandidateRepo, skillRepo *mocks.SkillRepo) {
				otherCandidate := domain.Candidate{ID: userId, UserId: uuid.New()}
				candidateRepo.On("GetByUserID", ctx, userId).Return(otherCandidate, nil).Once()
				repo.On("GetById", ctx, resumeID).Return(existing, nil).Once()
			},
			expectedError: domain.ErrForbidden,
		},
		{name: "candidate not found",
			req: Request{ID: resumeID, UserId: userId},
			mockSetup: func(repo *mocks.ResumeRepo, candidateRepo *mocks.CandidateRepo, skillRepo *mocks.SkillRepo) {
				candidateRepo.On("GetByUserID", ctx, userId).Return(domain.Candidate{}, domain.ErrCandidateNotFound).Once()
			},
			expectedError: domain.ErrCandidateNotFound,
		},
		{
			name: "skills not exist",
			req:  Request{ID: resumeID, UserId: userId, Data: &ResumeData{SkillsList: &[]uuid.UUID{uuid.New()}}},
			mockSetup: func(repo *mocks.ResumeRepo, candidateRepo *mocks.CandidateRepo, skillRepo *mocks.SkillRepo) {
				candidateRepo.On("GetByUserID", ctx, userId).Return(candidate, nil).Once()

				repo.On("GetById", ctx, resumeID).Return(existing, nil).Once()
				skillRepo.On("AreSkillsExist", ctx, mock.Anything).Return(false, nil).Once()
			},
			expectedError: domain.ErrSkillNotFound,
		},
		{
			name: "validation fail after update",
			req:  Request{ID: resumeID, UserId: userId, Data: &ResumeData{Phone: stringPtr("bad phone")}},
			mockSetup: func(repo *mocks.ResumeRepo, candidateRepo *mocks.CandidateRepo, skillRepo *mocks.SkillRepo) {
				candidateRepo.On("GetByUserID", ctx, userId).Return(candidate, nil).Once()

				repo.On("GetById", ctx, resumeID).Return(existing, nil).Once()
			},
			expectedError: domain.ErrInvalidPhoneFormat,
		},
		{
			name: "repo update error",
			req:  Request{ID: resumeID, UserId: userId, Name: stringPtr("New")},
			mockSetup: func(repo *mocks.ResumeRepo, candidateRepo *mocks.CandidateRepo, skillRepo *mocks.SkillRepo) {
				candidateRepo.On("GetByUserID", ctx, userId).Return(candidate, nil).Once()

				repo.On("GetById", ctx, resumeID).Return(existing, nil).Once()
				repo.On("Update", ctx, mock.AnythingOfType("*domain.Resume")).Return(ErrUpdateDb).Once()
			},
			expectedError: ErrUpdateDb,
		},
		{name: "candidate repo error",
			req: Request{ID: resumeID, UserId: userId},
			mockSetup: func(repo *mocks.ResumeRepo, candidateRepo *mocks.CandidateRepo, skillRepo *mocks.SkillRepo) {
				candidateRepo.On("GetByUserID", ctx, userId).Return(domain.Candidate{}, ErrCandidateDb).Once()
			},
			expectedError: ErrCandidateDb,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.ResumeRepo{}
			skillRepo := &mocks.SkillRepo{}
			candidateRepo := &mocks.CandidateRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(repo, candidateRepo, skillRepo)
			}

			uc := New(repo, skillRepo, candidateRepo)
			err := uc.Execute(ctx, tt.req)
			if tt.expectedError != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.expectedError), "expected %v got %v", tt.expectedError, err)
			} else {
				require.NoError(t, err)
			}

			repo.AssertExpectations(t)
			skillRepo.AssertExpectations(t)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
