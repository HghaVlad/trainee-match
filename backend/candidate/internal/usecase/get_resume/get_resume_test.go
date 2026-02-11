package get_resume

import (
	"context"
	"testing"

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

func (m *MockResumeRepo) GetByCandidateId(ctx context.Context, candidateId uuid.UUID) ([]domain.Resume, error) {
	args := m.Called(ctx, candidateId)
	resumes := args.Get(0)
	if resumes == nil {
		return []domain.Resume{}, args.Error(1)
	}
	return resumes.([]domain.Resume), args.Error(1)
}

func TestGetById(t *testing.T) {
	ctx := context.Background()

	tests := map[string]struct {
		req       GetByIdRequest
		setupMock func(*MockResumeRepo, uuid.UUID)
		wantErr   bool
		errCheck  func(error) bool
		check     func(*testing.T, *GetByIdResponse)
	}{
		"happy path: valid resume ID returns complete resume data": {
			req: GetByIdRequest{ID: uuid.New()},
			setupMock: func(m *MockResumeRepo, resumeID uuid.UUID) {
				candidateID := uuid.New()
				mockResume := domain.Resume{
					ID:          resumeID,
					CandidateId: candidateID,
					Name:        "Primary Resume",
					Status:      domain.Draft,
					Data: domain.ResumeData{
						FirstName:  "John",
						LastName:   "Doe",
						Email:      "john@example.com",
						SkillsList: []uuid.UUID{uuid.New()},
					},
				}
				m.On("GetById", ctx, resumeID).Return(mockResume, nil).Once()
			},
			wantErr: false,
			check: func(t *testing.T, resp *GetByIdResponse) {
				assert.NotNil(t, resp)
				assert.Equal(t, "Primary Resume", resp.Name)
				assert.NotZero(t, resp.ID)
			},
		},

		"not found: resume ID does not exist returns ErrResumeNotFound": {
			req: GetByIdRequest{ID: uuid.New()},
			setupMock: func(m *MockResumeRepo, resumeID uuid.UUID) {
				m.On("GetById", ctx, resumeID).
					Return(domain.Resume{}, domain.ErrResumeNotFound).Once()
			},
			wantErr: true,
			errCheck: func(err error) bool {
				return err == domain.ErrResumeNotFound
			},
		},

		"repository error: database connection failed": {
			req: GetByIdRequest{ID: uuid.New()},
			setupMock: func(m *MockResumeRepo, resumeID uuid.UUID) {
				m.On("GetById", ctx, resumeID).
					Return(domain.Resume{}, domain.ErrResumeNotFound).Once()
			},
			wantErr: true,
			errCheck: func(err error) bool {
				return err != nil
			},
		},

		"edge case: resume with minimal data": {
			req: GetByIdRequest{ID: uuid.New()},
			setupMock: func(m *MockResumeRepo, resumeID uuid.UUID) {
				mockResume := domain.Resume{
					ID:          resumeID,
					CandidateId: uuid.New(),
					Name:        "Basic",
					Status:      domain.Draft,
					Data: domain.ResumeData{
						FirstName: "Jane",
						LastName:  "Smith",
					},
				}
				m.On("GetById", ctx, resumeID).Return(mockResume, nil).Once()
			},
			wantErr: false,
			check: func(t *testing.T, resp *GetByIdResponse) {
				assert.Equal(t, "Basic", resp.Name)
			},
		},

		"edge case: nil UUID returns not found": {
			req: GetByIdRequest{ID: uuid.Nil},
			setupMock: func(m *MockResumeRepo, resumeID uuid.UUID) {
				m.On("GetById", ctx, resumeID).
					Return(domain.Resume{}, domain.ErrResumeNotFound).Once()
			},
			wantErr: true,
			errCheck: func(err error) bool {
				return err == domain.ErrResumeNotFound
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockRepo := new(MockResumeRepo)
			tt.setupMock(mockRepo, tt.req.ID)

			uc := New(mockRepo)
			result, err := uc.GetById(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				require.True(t, tt.errCheck(err), "error check failed: %v", err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.check != nil {
					tt.check(t, result)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetByCandidateId(t *testing.T) {
	ctx := context.Background()

	tests := map[string]struct {
		req       GetByCandidateIdRequest
		setupMock func(*MockResumeRepo, uuid.UUID)
		wantErr   bool
		errCheck  func(error) bool
		check     func(*testing.T, []*GetByCandidateIdResponse)
	}{
		"happy path: valid candidate ID returns multiple resumes": {
			req: GetByCandidateIdRequest{CandidateId: uuid.New()},
			setupMock: func(m *MockResumeRepo, candidateID uuid.UUID) {
				mockResumes := []domain.Resume{
					{
						ID:          uuid.New(),
						CandidateId: candidateID,
						Name:        "Primary",
						Status:      domain.Published,
					},
					{
						ID:          uuid.New(),
						CandidateId: candidateID,
						Name:        "Draft",
						Status:      domain.Draft,
					},
				}
				m.On("GetByCandidateId", ctx, candidateID).Return(mockResumes, nil).Once()
			},
			wantErr: false,
			check: func(t *testing.T, resumes []*GetByCandidateIdResponse) {
				assert.Len(t, resumes, 2)
				assert.Equal(t, "Primary", resumes[0].Name)
				assert.Equal(t, "Draft", resumes[1].Name)
			},
		},

		"not found: candidate has no resumes returns empty list": {
			req: GetByCandidateIdRequest{CandidateId: uuid.New()},
			setupMock: func(m *MockResumeRepo, candidateID uuid.UUID) {
				m.On("GetByCandidateId", ctx, candidateID).
					Return([]domain.Resume{}, nil).Once()
			},
			wantErr: false,
			check: func(t *testing.T, resumes []*GetByCandidateIdResponse) {
				assert.Len(t, resumes, 0)
			},
		},

		"repository error: database connection failed": {
			req: GetByCandidateIdRequest{CandidateId: uuid.New()},
			setupMock: func(m *MockResumeRepo, candidateID uuid.UUID) {
				m.On("GetByCandidateId", ctx, candidateID).
					Return(([]domain.Resume)(nil), domain.ErrResumeNotFound).Once()
			},
			wantErr: true,
			errCheck: func(err error) bool {
				return err != nil
			},
		},

		"edge case: candidate has single resume": {
			req: GetByCandidateIdRequest{CandidateId: uuid.New()},
			setupMock: func(m *MockResumeRepo, candidateID uuid.UUID) {
				mockResumes := []domain.Resume{
					{
						ID:          uuid.New(),
						CandidateId: candidateID,
						Name:        "Only Resume",
						Status:      domain.Published,
					},
				}
				m.On("GetByCandidateId", ctx, candidateID).Return(mockResumes, nil).Once()
			},
			wantErr: false,
			check: func(t *testing.T, resumes []*GetByCandidateIdResponse) {
				assert.Len(t, resumes, 1)
				assert.Equal(t, "Only Resume", resumes[0].Name)
			},
		},

		"edge case: candidate has many resumes": {
			req: GetByCandidateIdRequest{CandidateId: uuid.New()},
			setupMock: func(m *MockResumeRepo, candidateID uuid.UUID) {
				var mockResumes []domain.Resume
				for i := 0; i < 10; i++ {
					mockResumes = append(mockResumes, domain.Resume{
						ID:          uuid.New(),
						CandidateId: candidateID,
						Name:        "Resume",
						Status:      domain.Draft,
					})
				}
				m.On("GetByCandidateId", ctx, candidateID).Return(mockResumes, nil).Once()
			},
			wantErr: false,
			check: func(t *testing.T, resumes []*GetByCandidateIdResponse) {
				assert.Len(t, resumes, 10)
			},
		},

		"edge case: nil candidate ID returns empty or error": {
			req: GetByCandidateIdRequest{CandidateId: uuid.Nil},
			setupMock: func(m *MockResumeRepo, candidateID uuid.UUID) {
				m.On("GetByCandidateId", ctx, candidateID).
					Return([]domain.Resume{}, nil).Once()
			},
			wantErr: false,
			check: func(t *testing.T, resumes []*GetByCandidateIdResponse) {
				assert.Len(t, resumes, 0)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockRepo := new(MockResumeRepo)
			tt.setupMock(mockRepo, tt.req.CandidateId)

			uc := New(mockRepo)
			result, err := uc.GetByCandidateId(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				require.True(t, tt.errCheck(err), "error check failed: %v", err)
			} else {
				require.NoError(t, err)
				if tt.check != nil {
					tt.check(t, result)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
