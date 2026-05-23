package create_resume

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
)

func TestExecute(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	validReq := &Request{
		UserId: userID,
		Name:   "My Resume",
		Status: "draft",
		Data: ResumeData{
			LastName:        "Doe",
			FirstName:       "John",
			DateOfBirth:     time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
			Email:           "john@example.com",
			Phone:           "+1234567890",
			City:            "City",
			Citizenship:     "Country",
			Education:       []Education{{Level: "BSc", University: "Uni", StartYear: 2008, EndYear: 2012}},
			WorkExperiences: []WorkExperience{{Position: "Dev", Company: "Co", Period: "2012-2018"}},
			SkillsList:      []uuid.UUID{uuid.New()},
		},
	}
	candidate := domain.Candidate{ID: uuid.New(), UserId: userID}

	validId := uuid.New()
	var (
		errSkillDB  = errors.New("skill db error")
		errCreateDB = errors.New("create db error")
	)

	tests := []struct {
		name          string
		req           *Request
		mockSetup     func(*MockResumeRepo, *MockCandidateRepo, *MockSkillRepo)
		expectedID    uuid.UUID
		expectedError error
	}{
		{
			name: "valid request",
			req:  validReq,
			mockSetup: func(r *MockResumeRepo, c *MockCandidateRepo, s *MockSkillRepo) {
				c.On("GetByUserID", ctx, userID).Return(candidate, nil).Once()
				s.On("AreSkillsExist", ctx, mock.Anything).Return(true, nil).Once()
				r.On("Create", ctx, mock.AnythingOfType("*domain.Resume")).Return(validId, nil).Once()
			},
			expectedID:    validId,
			expectedError: nil,
		},
		{
			name: "missing email",
			req:  func() *Request { r := *validReq; r.Data.Email = ""; r.Data.SkillsList = nil; return &r }(),
			mockSetup: func(r *MockResumeRepo, c *MockCandidateRepo, s *MockSkillRepo) {
				c.On("GetByUserID", ctx, userID).Return(candidate, nil).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrInvalidEmailFormat,
		},
		{
			name: "missing phone",
			req:  func() *Request { r := *validReq; r.Data.Phone = ""; r.Data.SkillsList = nil; return &r }(),
			mockSetup: func(r *MockResumeRepo, c *MockCandidateRepo, s *MockSkillRepo) {
				c.On("GetByUserID", ctx, userID).Return(candidate, nil).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrInvalidPhoneFormat,
		},
		{
			name: "missing city",
			req:  func() *Request { r := *validReq; r.Data.City = ""; r.Data.SkillsList = nil; return &r }(),
			mockSetup: func(r *MockResumeRepo, c *MockCandidateRepo, s *MockSkillRepo) {
				c.On("GetByUserID", ctx, userID).Return(candidate, nil).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrInvalidCityFormat,
		},
		{
			name: "missing citizenship",
			req:  func() *Request { r := *validReq; r.Data.Citizenship = ""; r.Data.SkillsList = nil; return &r }(),
			mockSetup: func(r *MockResumeRepo, c *MockCandidateRepo, s *MockSkillRepo) {
				c.On("GetByUserID", ctx, userID).Return(candidate, nil).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrInvalidCitizenship,
		},
		{
			name: "empty education",
			req:  func() *Request { r := *validReq; r.Data.Education = nil; r.Data.SkillsList = nil; return &r }(),
			mockSetup: func(r *MockResumeRepo, c *MockCandidateRepo, s *MockSkillRepo) {
				c.On("GetByUserID", ctx, userID).Return(candidate, nil).Once()
				r.On("Create", ctx, mock.AnythingOfType("*domain.Resume")).Return(validId, nil).Once()
			},
			expectedID:    validId,
			expectedError: nil,
		},
		{
			name: "empty work experiences",
			req:  func() *Request { r := *validReq; r.Data.WorkExperiences = nil; r.Data.SkillsList = nil; return &r }(),
			mockSetup: func(r *MockResumeRepo, c *MockCandidateRepo, s *MockSkillRepo) {
				c.On("GetByUserID", ctx, userID).Return(candidate, nil).Once()
				r.On("Create", ctx, mock.AnythingOfType("*domain.Resume")).Return(validId, nil).Once()
			},
			expectedID:    validId,
			expectedError: nil,
		},
		{
			name: "skills repo error",
			req:  func() *Request { r := *validReq; r.Data.SkillsList = []uuid.UUID{uuid.New()}; return &r }(),
			mockSetup: func(r *MockResumeRepo, c *MockCandidateRepo, s *MockSkillRepo) {
				c.On("GetByUserID", ctx, userID).Return(candidate, nil).Once()
				s.On("AreSkillsExist", ctx, mock.Anything).Return(false, errSkillDB).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: errSkillDB,
		},
		{
			name: "skills not exist",
			req:  func() *Request { r := *validReq; r.Data.SkillsList = []uuid.UUID{uuid.New()}; return &r }(),
			mockSetup: func(r *MockResumeRepo, c *MockCandidateRepo, s *MockSkillRepo) {
				c.On("GetByUserID", ctx, userID).Return(candidate, nil).Once()

				s.On("AreSkillsExist", ctx, mock.Anything).Return(false, nil).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrSkillNotFound,
		},
		{
			name: "repo create error",
			req:  validReq,
			mockSetup: func(r *MockResumeRepo, c *MockCandidateRepo, s *MockSkillRepo) {
				c.On("GetByUserID", ctx, userID).Return(candidate, nil).Once()
				s.On("AreSkillsExist", ctx, mock.Anything).Return(true, nil).Once()
				r.On("Create", ctx, mock.AnythingOfType("*domain.Resume")).Return(uuid.Nil, errCreateDB).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: errCreateDB,
		},
		{
			name: "candidate not found",
			req:  validReq,
			mockSetup: func(r *MockResumeRepo, c *MockCandidateRepo, s *MockSkillRepo) {
				c.On("GetByUserID", ctx, userID).Return(domain.Candidate{}, domain.ErrCandidateNotFound).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrCandidateNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockResumeRepo{}
			skillRepo := &MockSkillRepo{}
			candidateRepo := &MockCandidateRepo{}
			if tt.mockSetup != nil {
				tt.mockSetup(repo, candidateRepo, skillRepo)
			}
			writer := &MockEventWriter{}
			writer.On("WriteResumeUpserted", context.Background(), mock.AnythingOfType("events.ResumeUpserted")).
				Return(nil)

			trManager := &MockTrManager{}
			trManager.On("Do", mock.Anything, mock.Anything).
				Return(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})

			uc := New(repo, skillRepo, candidateRepo, writer, trManager)
			res, err := uc.Execute(ctx, *tt.req)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
				require.Equal(t, uuid.Nil, res.ID)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedID, res.ID)
			}

			repo.AssertExpectations(t)
			skillRepo.AssertExpectations(t)
		})
	}
}
