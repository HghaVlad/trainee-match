package create_resume

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

func (m *MockResumeRepo) Create(ctx context.Context, resume *domain.Resume) (uuid.UUID, error) {
	args := m.Called(ctx, resume)
	var id uuid.UUID
	if args.Get(0) != nil {
		id = args.Get(0).(uuid.UUID)
	}
	return id, args.Error(1)
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
	candidateID := uuid.New()

	tests := map[string]struct {
		req         Request
		setupMock   func(*MockResumeRepo, *MockSkillRepo, *Request)
		wantErr     bool
		errCheck    func(error) bool
		resultCheck func(*testing.T, Response)
	}{
		"happy path: valid resume creation with all fields": {
			req: Request{
				CandidateId: candidateID,
				Name:        "Primary Resume",
				Status:      domain.Draft,
				Data: ResumeData{
					LastName:      "Doe",
					FirstName:     "John",
					MiddleName:    "Alexander",
					DateOfBirth:   time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
					Email:         "john@example.com",
					Phone:         "+1234567890",
					City:          "New York",
					Citizenship:   "USA",
					SkillsList:    []uuid.UUID{uuid.New(), uuid.New()},
					EnglishLevel:  "Fluent",
					DesiredFormat: "PDF",
				},
			},
			setupMock: func(resumeRepo *MockResumeRepo, skillRepo *MockSkillRepo, req *Request) {
				// Mock skill validation
				skillRepo.On("AreSkillsExist", ctx, mock.MatchedBy(func(ids []uuid.UUID) bool {
					return len(ids) == 2
				})).Return(true, nil).Once()

				// Mock resume creation
				resumeID := uuid.New()
				resumeRepo.On("Create", ctx, mock.AnythingOfType("*domain.Resume")).
					Return(resumeID, nil).Once()
			},
			wantErr: false,
			resultCheck: func(t *testing.T, resp Response) {
				assert.NotEqual(t, uuid.Nil, resp.ID)
			},
		},

		"not found: skills do not exist returns ErrSkillNotFound": {
			req: Request{
				CandidateId: candidateID,
				Name:        "Invalid Resume",
				Status:      domain.Published,
				Data: ResumeData{
					FirstName:  "John",
					LastName:   "Doe",
					SkillsList: []uuid.UUID{uuid.New()},
				},
			},
			setupMock: func(resumeRepo *MockResumeRepo, skillRepo *MockSkillRepo, req *Request) {
				skillRepo.On("AreSkillsExist", ctx, mock.AnythingOfType("[]uuid.UUID")).
					Return(false, nil).Once()
			},
			wantErr: true,
			errCheck: func(err error) bool {
				return err == domain.ErrSkillNotFound
			},
		},

		"repository error: database connection failed on resume create": {
			req: Request{
				CandidateId: candidateID,
				Name:        "Test Resume",
				Status:      domain.Draft,
				Data: ResumeData{
					FirstName:  "Jane",
					LastName:   "Smith",
					SkillsList: []uuid.UUID{uuid.New()},
				},
			},
			setupMock: func(resumeRepo *MockResumeRepo, skillRepo *MockSkillRepo, req *Request) {
				skillRepo.On("AreSkillsExist", ctx, mock.AnythingOfType("[]uuid.UUID")).
					Return(true, nil).Once()

				resumeRepo.On("Create", ctx, mock.AnythingOfType("*domain.Resume")).
					Return(uuid.Nil, domain.ErrResumeNotFound).Once()
			},
			wantErr: true,
			errCheck: func(err error) bool {
				return err != nil
			},
		},

		"edge case: resume with empty skills list": {
			req: Request{
				CandidateId: candidateID,
				Name:        "No Skills Resume",
				Status:      domain.Draft,
				Data: ResumeData{
					FirstName:  "Bob",
					LastName:   "Johnson",
					SkillsList: []uuid.UUID{},
				},
			},
			setupMock: func(resumeRepo *MockResumeRepo, skillRepo *MockSkillRepo, req *Request) {
				skillRepo.On("AreSkillsExist", ctx, mock.MatchedBy(func(ids []uuid.UUID) bool {
					return len(ids) == 0
				})).Return(true, nil).Once()

				resumeID := uuid.New()
				resumeRepo.On("Create", ctx, mock.AnythingOfType("*domain.Resume")).
					Return(resumeID, nil).Once()
			},
			wantErr: false,
			resultCheck: func(t *testing.T, resp Response) {
				assert.NotEqual(t, uuid.Nil, resp.ID)
			},
		},

		"edge case: resume with only required fields": {
			req: Request{
				CandidateId: candidateID,
				Name:        "Minimal Resume",
				Status:      domain.Draft,
				Data: ResumeData{
					FirstName:  "Alice",
					LastName:   "Williams",
					SkillsList: []uuid.UUID{},
				},
			},
			setupMock: func(resumeRepo *MockResumeRepo, skillRepo *MockSkillRepo, req *Request) {
				skillRepo.On("AreSkillsExist", ctx, mock.MatchedBy(func(ids []uuid.UUID) bool {
					return len(ids) == 0
				})).Return(true, nil).Once()

				resumeID := uuid.New()
				resumeRepo.On("Create", ctx, mock.AnythingOfType("*domain.Resume")).
					Return(resumeID, nil).Once()
			},
			wantErr: false,
			resultCheck: func(t *testing.T, resp Response) {
				assert.NotEqual(t, uuid.Nil, resp.ID)
			},
		},

		"edge case: resume with full education and experience history": {
			req: Request{
				CandidateId: candidateID,
				Name:        "Full Resume",
				Status:      domain.Published,
				Data: ResumeData{
					FirstName: "David",
					LastName:  "Brown",
					Email:     "david@example.com",
					Phone:     "+9876543210",
					City:      "Boston",
					Education: []Education{
						{
							Level:          "Bachelor",
							University:     "MIT",
							Faculty:        "Engineering",
							Specialization: "Computer Science",
							StartYear:      2016,
							EndYear:        2020,
							Format:         "Full-time",
						},
					},
					WorkExperiences: []WorkExperience{
						{
							Position:         "Senior Engineer",
							Company:          "Tech Corp",
							Period:           "2020-2024",
							Responsibilities: "Led development team",
						},
					},
					SkillsList: []uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
				},
			},
			setupMock: func(resumeRepo *MockResumeRepo, skillRepo *MockSkillRepo, req *Request) {
				skillRepo.On("AreSkillsExist", ctx, mock.MatchedBy(func(ids []uuid.UUID) bool {
					return len(ids) == 3
				})).Return(true, nil).Once()

				resumeID := uuid.New()
				resumeRepo.On("Create", ctx, mock.AnythingOfType("*domain.Resume")).
					Return(resumeID, nil).Once()
			},
			wantErr: false,
			resultCheck: func(t *testing.T, resp Response) {
				assert.NotEqual(t, uuid.Nil, resp.ID)
			},
		},

		"edge case: skill repository error returns wrapped error": {
			req: Request{
				CandidateId: candidateID,
				Name:        "Resume with Skill Error",
				Status:      domain.Draft,
				Data: ResumeData{
					FirstName:  "Eve",
					LastName:   "Davis",
					SkillsList: []uuid.UUID{uuid.New()},
				},
			},
			setupMock: func(resumeRepo *MockResumeRepo, skillRepo *MockSkillRepo, req *Request) {
				skillRepo.On("AreSkillsExist", ctx, mock.AnythingOfType("[]uuid.UUID")).
					Return(false, domain.ErrSkillNotFound).Once()
			},
			wantErr: true,
			errCheck: func(err error) bool {
				return err != nil
			},
		},

		"edge case: nil candidate ID": {
			req: Request{
				CandidateId: uuid.Nil,
				Name:        "Resume",
				Status:      domain.Draft,
				Data: ResumeData{
					FirstName:  "Frank",
					LastName:   "Miller",
					SkillsList: []uuid.UUID{},
				},
			},
			setupMock: func(resumeRepo *MockResumeRepo, skillRepo *MockSkillRepo, req *Request) {
				skillRepo.On("AreSkillsExist", ctx, mock.MatchedBy(func(ids []uuid.UUID) bool {
					return len(ids) == 0
				})).Return(true, nil).Once()

				resumeID := uuid.New()
				resumeRepo.On("Create", ctx, mock.AnythingOfType("*domain.Resume")).
					Return(resumeID, nil).Once()
			},
			wantErr: false,
			resultCheck: func(t *testing.T, resp Response) {
				assert.NotEqual(t, uuid.Nil, resp.ID)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup
			mockResumeRepo := new(MockResumeRepo)
			mockSkillRepo := new(MockSkillRepo)
			tt.setupMock(mockResumeRepo, mockSkillRepo, &tt.req)

			uc := New(mockResumeRepo, mockSkillRepo)

			// Execute
			result, err := uc.Execute(ctx, tt.req)

			// Verify
			if tt.wantErr {
				require.Error(t, err)
				require.True(t, tt.errCheck(err), "error check failed: %v", err)
			} else {
				require.NoError(t, err)
				if tt.resultCheck != nil {
					tt.resultCheck(t, result)
				}
			}

			// Cleanup: Verify mock expectations
			mockResumeRepo.AssertExpectations(t)
			mockSkillRepo.AssertExpectations(t)
		})
	}
}
