package get_skill

import (
	"context"
	"testing"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockSkillRepo implements the SkillRepo interface
type MockSkillRepo struct {
	mock.Mock
}

func (m *MockSkillRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.Skill, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Skill), args.Error(1)
}

func (m *MockSkillRepo) List(ctx context.Context) ([]domain.Skill, error) {
	args := m.Called(ctx)
	skills := args.Get(0)
	if skills == nil {
		return []domain.Skill{}, args.Error(1)
	}
	return skills.([]domain.Skill), args.Error(1)
}

func TestExecute(t *testing.T) {
	ctx := context.Background()

	tests := map[string]struct {
		req       GetByIdRequest
		setupMock func(*MockSkillRepo, uuid.UUID)
		wantErr   bool
		errCheck  func(error) bool
	}{
		"happy path: valid skill ID returns skill data": {
			req: GetByIdRequest{ID: uuid.New()},
			setupMock: func(m *MockSkillRepo, skillID uuid.UUID) {
				m.On("GetByID", ctx, skillID).
					Return(domain.Skill{ID: skillID, Name: "Go Programming"}, nil).Once()
			},
			wantErr: false,
		},

		"not found: skill ID does not exist returns ErrSkillNotFound": {
			req: GetByIdRequest{ID: uuid.New()},
			setupMock: func(m *MockSkillRepo, skillID uuid.UUID) {
				m.On("GetByID", ctx, skillID).
					Return(domain.Skill{}, domain.ErrSkillNotFound).Once()
			},
			wantErr: true,
			errCheck: func(err error) bool {
				return err == domain.ErrSkillNotFound
			},
		},

		"repository error: database connection failed": {
			req: GetByIdRequest{ID: uuid.New()},
			setupMock: func(m *MockSkillRepo, skillID uuid.UUID) {
				m.On("GetByID", ctx, skillID).
					Return(domain.Skill{}, domain.ErrSkillNotFound).Once()
			},
			wantErr: true,
			errCheck: func(err error) bool {
				return err != nil
			},
		},

		"edge case: skill with empty name": {
			req: GetByIdRequest{ID: uuid.New()},
			setupMock: func(m *MockSkillRepo, skillID uuid.UUID) {
				m.On("GetByID", ctx, skillID).
					Return(domain.Skill{ID: skillID, Name: ""}, nil).Once()
			},
			wantErr: false,
		},

		"edge case: skill with special characters (C++/Python/машинное обучение)": {
			req: GetByIdRequest{ID: uuid.New()},
			setupMock: func(m *MockSkillRepo, skillID uuid.UUID) {
				m.On("GetByID", ctx, skillID).
					Return(domain.Skill{ID: skillID, Name: "C++/Python-машинное обучение"}, nil).Once()
			},
			wantErr: false,
		},

		"edge case: nil UUID returns not found": {
			req: GetByIdRequest{ID: uuid.Nil},
			setupMock: func(m *MockSkillRepo, skillID uuid.UUID) {
				m.On("GetByID", ctx, skillID).
					Return(domain.Skill{}, domain.ErrSkillNotFound).Once()
			},
			wantErr: true,
			errCheck: func(err error) bool {
				return err == domain.ErrSkillNotFound
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockSkillRepo)
			tt.setupMock(mockRepo, tt.req.ID)

			uc := New(mockRepo)

			// Execute
			result, err := uc.Execute(ctx, tt.req)

			// Verify
			if tt.wantErr {
				require.Error(t, err)
				if tt.errCheck != nil {
					require.True(t, tt.errCheck(err), "error check failed: %v", err)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.NotZero(t, result.ID)
			}

			// Cleanup: Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestExecuteList(t *testing.T) {
	ctx := context.Background()

	tests := map[string]struct {
		req       ListRequest
		setupMock func(*MockSkillRepo)
		wantErr   bool
		errCheck  func(error) bool
		checkLen  func(*testing.T, []*ListResponse)
	}{
		"happy path: list all skills returns multiple skills": {
			req: ListRequest{},
			setupMock: func(m *MockSkillRepo) {
				skills := []domain.Skill{
					{ID: uuid.New(), Name: "Go"},
					{ID: uuid.New(), Name: "Python"},
					{ID: uuid.New(), Name: "Rust"},
				}
				m.On("List", ctx).Return(skills, nil).Once()
			},
			wantErr: false,
			checkLen: func(t *testing.T, skills []*ListResponse) {
				assert.Len(t, skills, 3)
				assert.Equal(t, "Go", skills[0].Name)
				assert.Equal(t, "Python", skills[1].Name)
				assert.Equal(t, "Rust", skills[2].Name)
			},
		},

		"not found: no skills exist returns empty list": {
			req: ListRequest{},
			setupMock: func(m *MockSkillRepo) {
				m.On("List", ctx).Return([]domain.Skill{}, nil).Once()
			},
			wantErr: false,
			checkLen: func(t *testing.T, skills []*ListResponse) {
				assert.Len(t, skills, 0)
			},
		},

		"repository error: database connection failed": {
			req: ListRequest{},
			setupMock: func(m *MockSkillRepo) {
				m.On("List", ctx).
					Return(([]domain.Skill)(nil), domain.ErrSkillNotFound).Once()
			},
			wantErr: true,
			errCheck: func(err error) bool {
				return err != nil
			},
		},

		"edge case: list returns single skill": {
			req: ListRequest{},
			setupMock: func(m *MockSkillRepo) {
				skills := []domain.Skill{
					{ID: uuid.New(), Name: "JavaScript"},
				}
				m.On("List", ctx).Return(skills, nil).Once()
			},
			wantErr: false,
			checkLen: func(t *testing.T, skills []*ListResponse) {
				assert.Len(t, skills, 1)
				assert.Equal(t, "JavaScript", skills[0].Name)
			},
		},

		"edge case: list returns many skills (stress test with 100 items)": {
			req: ListRequest{},
			setupMock: func(m *MockSkillRepo) {
				var skills []domain.Skill
				for i := 0; i < 100; i++ {
					skills = append(skills, domain.Skill{ID: uuid.New(), Name: "Skill"})
				}
				m.On("List", ctx).Return(skills, nil).Once()
			},
			wantErr: false,
			checkLen: func(t *testing.T, skills []*ListResponse) {
				assert.Len(t, skills, 100)
			},
		},

		"edge case: list includes skills with special characters": {
			req: ListRequest{},
			setupMock: func(m *MockSkillRepo) {
				skills := []domain.Skill{
					{ID: uuid.New(), Name: "C++"},
					{ID: uuid.New(), Name: "Node.js"},
					{ID: uuid.New(), Name: "Machine Learning/AI"},
					{ID: uuid.New(), Name: "Web-разработка"},
				}
				m.On("List", ctx).Return(skills, nil).Once()
			},
			wantErr: false,
			checkLen: func(t *testing.T, skills []*ListResponse) {
				assert.Len(t, skills, 4)
				assert.Equal(t, "C++", skills[0].Name)
				assert.Equal(t, "Web-разработка", skills[3].Name)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockSkillRepo)
			tt.setupMock(mockRepo)

			uc := New(mockRepo)

			// Execute
			result, err := uc.ExecuteList(ctx, tt.req)

			// Verify
			if tt.wantErr {
				require.Error(t, err)
				if tt.errCheck != nil {
					require.True(t, tt.errCheck(err), "error check failed: %v", err)
				}
			} else {
				require.NoError(t, err)
				if tt.checkLen != nil {
					tt.checkLen(t, result)
				}
			}

			// Cleanup: Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}
