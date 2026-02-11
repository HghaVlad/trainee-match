package get_candidate_by_user_id

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

// MockCandidateRepo is a mock implementation of the CandidateRepo interface
type MockCandidateRepo struct {
	mock.Mock
}

func (m *MockCandidateRepo) GetByUserID(ctx context.Context, id uuid.UUID) (domain.Candidate, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Candidate), args.Error(1)
}

func TestExecute(t *testing.T) {
	ctx := context.Background()

	tests := map[string]struct {
		setupMock  func(*MockCandidateRepo, uuid.UUID)
		inputID    uuid.UUID
		wantResult *CandidateResponse
		wantErr    bool
		errCheck   func(error) bool
	}{
		"happy path: valid user ID returns candidate data": {
			inputID: uuid.New(),
			setupMock: func(m *MockCandidateRepo, userID uuid.UUID) {
				candidateID := uuid.New()
				birthday := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
				expectedCandidate := domain.Candidate{
					ID:       candidateID,
					UserId:   userID,
					Phone:    "+1234567890",
					Telegram: "@testuser",
					City:     "New York",
					Birthday: birthday,
				}
				m.On("GetByUserID", ctx, userID).Return(expectedCandidate, nil).Once()
			},
			wantErr: false,
			errCheck: func(err error) bool {
				return err == nil
			},
		},

		"not found: user ID does not have candidate returns ErrCandidateNotFound": {
			inputID: uuid.New(),
			setupMock: func(m *MockCandidateRepo, userID uuid.UUID) {
				m.On("GetByUserID", ctx, userID).
					Return(domain.Candidate{}, domain.ErrCandidateNotFound).
					Once()
			},
			wantErr: true,
			errCheck: func(err error) bool {
				return err == domain.ErrCandidateNotFound
			},
		},

		"repository error: database query failure returns error": {
			inputID: uuid.New(),
			setupMock: func(m *MockCandidateRepo, userID uuid.UUID) {
				m.On("GetByUserID", ctx, userID).
					Return(domain.Candidate{}, domain.ErrCandidateNotFound).
					Once()
			},
			wantErr: true,
			errCheck: func(err error) bool {
				return err != nil
			},
		},

		"edge case: candidate with minimal fields (no phone, telegram, city)": {
			inputID: uuid.New(),
			setupMock: func(m *MockCandidateRepo, userID uuid.UUID) {
				candidateID := uuid.New()
				birthday := time.Date(1985, 5, 20, 0, 0, 0, 0, time.UTC)
				expectedCandidate := domain.Candidate{
					ID:       candidateID,
					UserId:   userID,
					Phone:    "",
					Telegram: "",
					City:     "",
					Birthday: birthday,
				}
				m.On("GetByUserID", ctx, userID).Return(expectedCandidate, nil).Once()
			},
			wantErr: false,
			errCheck: func(err error) bool {
				return err == nil
			},
		},

		"edge case: candidate with international phone and cyrillic city": {
			inputID: uuid.New(),
			setupMock: func(m *MockCandidateRepo, userID uuid.UUID) {
				candidateID := uuid.New()
				birthday := time.Date(1992, 8, 10, 0, 0, 0, 0, time.UTC)
				expectedCandidate := domain.Candidate{
					ID:       candidateID,
					UserId:   userID,
					Phone:    "+7 (999) 999-99-99",
					Telegram: "@moscow_user_2024",
					City:     "Москва",
					Birthday: birthday,
				}
				m.On("GetByUserID", ctx, userID).Return(expectedCandidate, nil).Once()
			},
			wantErr: false,
			errCheck: func(err error) bool {
				return err == nil
			},
		},

		"edge case: nil UUID as user ID returns not found": {
			inputID: uuid.Nil,
			setupMock: func(m *MockCandidateRepo, userID uuid.UUID) {
				m.On("GetByUserID", ctx, userID).
					Return(domain.Candidate{}, domain.ErrCandidateNotFound).
					Once()
			},
			wantErr: true,
			errCheck: func(err error) bool {
				return err == domain.ErrCandidateNotFound
			},
		},

		"edge case: candidate with very long text fields": {
			inputID: uuid.New(),
			setupMock: func(m *MockCandidateRepo, userID uuid.UUID) {
				candidateID := uuid.New()
				birthday := time.Date(1988, 3, 25, 0, 0, 0, 0, time.UTC)
				expectedCandidate := domain.Candidate{
					ID:       candidateID,
					UserId:   userID,
					Phone:    "+1" + string(make([]byte, 50)),
					Telegram: "@" + string(make([]byte, 100)),
					City:     string(make([]byte, 150)),
					Birthday: birthday,
				}
				m.On("GetByUserID", ctx, userID).Return(expectedCandidate, nil).Once()
			},
			wantErr: false,
			errCheck: func(err error) bool {
				return err == nil
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockCandidateRepo)
			tt.setupMock(mockRepo, tt.inputID)

			uc := New(mockRepo)

			// Execute
			result, err := uc.Execute(ctx, tt.inputID)

			// Verify error expectation
			if tt.wantErr {
				require.Error(t, err)
				require.True(t, tt.errCheck(err), "error does not match expectation: %v", err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.NotZero(t, result.ID, "result ID should not be zero")
				assert.Equal(t, tt.inputID, result.UserID, "result UserID should match input")
			}

			// Cleanup: Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}
