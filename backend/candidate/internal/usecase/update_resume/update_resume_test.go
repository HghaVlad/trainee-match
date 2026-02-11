package update_resume

import (
	"context"
	"testing"
	"time"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockResumeRepo is a mock implementation of the ResumeRepo interface
type MockResumeRepo struct {
	mock.Mock
}

func (m *MockResumeRepo) GetById(ctx context.Context, id uuid.UUID) (domain.Resume, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Resume), args.Error(1)
}

func (m *MockResumeRepo) Update(ctx context.Context, resume *domain.Resume) error {
	args := m.Called(ctx, resume)
	return args.Error(0)
}

// MockSkillRepo is a mock implementation of the SkillRepo interface
type MockSkillRepo struct {
	mock.Mock
}

func (m *MockSkillRepo) AreSkillsExist(ctx context.Context, ids []uuid.UUID) (bool, error) {
	args := m.Called(ctx, ids)
	return args.Bool(0), args.Error(1)
}

func TestExecute(t *testing.T) {
	ctx := context.Background()
	resumeID := uuid.New()
	candidateID := uuid.New()

	tests := map[string]struct {
		req       *Request
		setupMock func(*MockResumeRepo, *MockSkillRepo, *Request)
		wantErr   bool
		errCheck  func(error) bool
	}{
		"happy path: update resume name successfully": {
			req: &Request{
				ID:   resumeID,
				Name: stringPtr("Updated Resume"),
			},
			setupMock: func(resumeRepo *MockResumeRepo, skillRepo *MockSkillRepo, req *Request) {
				existingResume := domain.Resume{
					ID:          resumeID,
					CandidateId: candidateID,
					Name:        "Old Name",
					Status:      domain.Draft,
					Data:        domain.ResumeData{},
				}
				resumeRepo.On("GetById", ctx, resumeID).Return(existingResume, nil).Once()
				resumeRepo.On("Update", ctx, mock.AnythingOfType("*domain.Resume")).Return(nil).Once()
			},
			wantErr: false,
		},

		"not found: resume ID does not exist returns ErrResumeNotFound": {
			req: &Request{
				ID: uuid.New(),
			},
			setupMock: func(resumeRepo *MockResumeRepo, skillRepo *MockSkillRepo, req *Request) {
				resumeRepo.On("GetById", ctx, req.ID).
					Return(domain.Resume{}, domain.ErrResumeNotFound).Once()
			},
			wantErr: true,
			errCheck: func(err error) bool {
				return err == domain.ErrResumeNotFound
			},
		},

		"repository error: database connection failed on update": {
			req: &Request{
				ID:     resumeID,
				Status: intPtr(domain.Published),
			},
			setupMock: func(resumeRepo *MockResumeRepo, skillRepo *MockSkillRepo, req *Request) {
				existingResume := domain.Resume{
					ID:          resumeID,
					CandidateId: candidateID,
					Name:        "Resume",
					Status:      domain.Draft,
					Data:        domain.ResumeData{},
				}
				resumeRepo.On("GetById", ctx, resumeID).Return(existingResume, nil).Once()
				resumeRepo.On("Update", ctx, mock.AnythingOfType("*domain.Resume")).
					Return(domain.ErrResumeNotFound).Once()
			},
			wantErr: true,
			errCheck: func(err error) bool {
				return err != nil
			},
		},

		"edge case: update all resume fields simultaneously": {
			req: &Request{
				ID:     resumeID,
				Name:   stringPtr("Fully Updated Resume"),
				Status: intPtr(domain.Published),
				Data: &ResumeData{
					FirstName:  stringPtr("John"),
					LastName:   stringPtr("Doe"),
					Email:      stringPtr("john@example.com"),
					Phone:      stringPtr("+1234567890"),
					City:       stringPtr("New York"),
					SkillsList: &[]uuid.UUID{uuid.New(), uuid.New()},
				},
			},
			setupMock: func(resumeRepo *MockResumeRepo, skillRepo *MockSkillRepo, req *Request) {
				existingResume := domain.Resume{
					ID:          resumeID,
					CandidateId: candidateID,
					Name:        "Old",
					Status:      domain.Draft,
					Data:        domain.ResumeData{},
				}
				resumeRepo.On("GetById", ctx, resumeID).Return(existingResume, nil).Once()
				skillRepo.On("AreSkillsExist", ctx, mock.MatchedBy(func(ids []uuid.UUID) bool {
					return len(ids) == 2
				})).Return(true, nil).Once()
				resumeRepo.On("Update", ctx, mock.AnythingOfType("*domain.Resume")).Return(nil).Once()
			},
			wantErr: false,
		},

		"edge case: update skills validates they exist": {
			req: &Request{
				ID: resumeID,
				Data: &ResumeData{
					SkillsList: &[]uuid.UUID{uuid.New()},
				},
			},
			setupMock: func(resumeRepo *MockResumeRepo, skillRepo *MockSkillRepo, req *Request) {
				existingResume := domain.Resume{
					ID:          resumeID,
					CandidateId: candidateID,
					Name:        "Resume",
					Status:      domain.Draft,
					Data:        domain.ResumeData{},
				}
				resumeRepo.On("GetById", ctx, resumeID).Return(existingResume, nil).Once()
				skillRepo.On("AreSkillsExist", ctx, mock.AnythingOfType("[]uuid.UUID")).
					Return(true, nil).Once()
				resumeRepo.On("Update", ctx, mock.AnythingOfType("*domain.Resume")).Return(nil).Once()
			},
			wantErr: false,
		},

		"edge case: update with non-existent skills returns ErrSkillNotFound": {
			req: &Request{
				ID: resumeID,
				Data: &ResumeData{
					SkillsList: &[]uuid.UUID{uuid.New()},
				},
			},
			setupMock: func(resumeRepo *MockResumeRepo, skillRepo *MockSkillRepo, req *Request) {
				existingResume := domain.Resume{
					ID:          resumeID,
					CandidateId: candidateID,
					Name:        "Resume",
					Status:      domain.Draft,
					Data:        domain.ResumeData{},
				}
				resumeRepo.On("GetById", ctx, resumeID).Return(existingResume, nil).Once()
				skillRepo.On("AreSkillsExist", ctx, mock.AnythingOfType("[]uuid.UUID")).
					Return(false, nil).Once()
			},
			wantErr: true,
			errCheck: func(err error) bool {
				return err == domain.ErrSkillNotFound
			},
		},

		"edge case: update with empty skills list": {
			req: &Request{
				ID: resumeID,
				Data: &ResumeData{
					SkillsList: &[]uuid.UUID{},
				},
			},
			setupMock: func(resumeRepo *MockResumeRepo, skillRepo *MockSkillRepo, req *Request) {
				existingResume := domain.Resume{
					ID:          resumeID,
					CandidateId: candidateID,
					Name:        "Resume",
					Status:      domain.Draft,
					Data:        domain.ResumeData{},
				}
				resumeRepo.On("GetById", ctx, resumeID).Return(existingResume, nil).Once()
				skillRepo.On("AreSkillsExist", ctx, mock.MatchedBy(func(ids []uuid.UUID) bool {
					return len(ids) == 0
				})).Return(true, nil).Once()
				resumeRepo.On("Update", ctx, mock.AnythingOfType("*domain.Resume")).Return(nil).Once()
			},
			wantErr: false,
		},

		"edge case: update with no fields specified": {
			req: &Request{
				ID: resumeID,
			},
			setupMock: func(resumeRepo *MockResumeRepo, skillRepo *MockSkillRepo, req *Request) {
				existingResume := domain.Resume{
					ID:          resumeID,
					CandidateId: candidateID,
					Name:        "Resume",
					Status:      domain.Draft,
					Data:        domain.ResumeData{},
				}
				resumeRepo.On("GetById", ctx, resumeID).Return(existingResume, nil).Once()
				resumeRepo.On("Update", ctx, mock.AnythingOfType("*domain.Resume")).Return(nil).Once()
			},
			wantErr: false,
		},

		"edge case: nil resume ID returns not found": {
			req: &Request{
				ID: uuid.Nil,
			},
			setupMock: func(resumeRepo *MockResumeRepo, skillRepo *MockSkillRepo, req *Request) {
				resumeRepo.On("GetById", ctx, uuid.Nil).
					Return(domain.Resume{}, domain.ErrResumeNotFound).Once()
			},
			wantErr: true,
			errCheck: func(err error) bool {
				return err == domain.ErrResumeNotFound
			},
		},

		"edge case: update with international text and special characters": {
			req: &Request{
				ID:   resumeID,
				Name: stringPtr("Резюме-2024-Final_v3"),
				Data: &ResumeData{
					FirstName: stringPtr("Jean-Pierre"),
					LastName:  stringPtr("Müller"),
					City:      stringPtr("München, Bavière, Allemagne"),
					Phone:     stringPtr("+49 (89) 123-456-78"),
				},
			},
			setupMock: func(resumeRepo *MockResumeRepo, skillRepo *MockSkillRepo, req *Request) {
				existingResume := domain.Resume{
					ID:          resumeID,
					CandidateId: candidateID,
					Name:        "Resume",
					Status:      domain.Draft,
					Data:        domain.ResumeData{},
				}
				resumeRepo.On("GetById", ctx, resumeID).Return(existingResume, nil).Once()
				resumeRepo.On("Update", ctx, mock.AnythingOfType("*domain.Resume")).Return(nil).Once()
			},
			wantErr: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup
			mockResumeRepo := new(MockResumeRepo)
			mockSkillRepo := new(MockSkillRepo)
			tt.setupMock(mockResumeRepo, mockSkillRepo, tt.req)

			uc := New(mockResumeRepo, mockSkillRepo)

			// Execute
			result, err := uc.Execute(ctx, *tt.req)

			// Verify
			if tt.wantErr {
				require.Error(t, err)
				if tt.errCheck != nil {
					require.True(t, tt.errCheck(err), "error check failed: %v", err)
				}
			} else {
				require.NoError(t, err)
				assert.True(t, result.Success)
			}

			// Cleanup: Verify mock expectations
			mockResumeRepo.AssertExpectations(t)
			mockSkillRepo.AssertExpectations(t)
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func timePtr(t time.Time) *time.Time {
	return &t
}
