package update_candidate

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

func (m *MockCandidateRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.Candidate, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Candidate), args.Error(1)
}

func (m *MockCandidateRepo) Update(ctx context.Context, candidate domain.Candidate) (domain.Candidate, error) {
	args := m.Called(ctx, candidate)
	return args.Get(0).(domain.Candidate), args.Error(1)
}

func TestExecute(t *testing.T) {
	ctx := context.Background()
	candidateID := uuid.New()
	userID := uuid.New()
	birthday := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		req         *Request
		setupMock   func(*MockCandidateRepo, *Request)
		wantErr     bool
		errCheck    func(error) bool
		resultCheck func(*testing.T, *CandidateResponse)
	}{
		"happy path: update one field successfully": {
			req: &Request{
				ID:    candidateID,
				Phone: stringPtr("+9999999999"),
			},
			setupMock: func(m *MockCandidateRepo, req *Request) {
				// GetByID returns existing candidate
				existingCandidate := domain.Candidate{
					ID:       candidateID,
					UserId:   userID,
					Phone:    "+1111111111",
					Telegram: "@olduser",
					City:     "Old City",
					Birthday: birthday,
				}
				m.On("GetByID", ctx, req.ID).Return(existingCandidate, nil).Once()

				// Update returns updated candidate
				updatedCandidate := domain.Candidate{
					ID:       candidateID,
					UserId:   userID,
					Phone:    "+9999999999",
					Telegram: "@olduser",
					City:     "Old City",
					Birthday: birthday,
				}
				m.On("Update", ctx, mock.MatchedBy(func(c domain.Candidate) bool {
					return c.Phone == "+9999999999"
				})).Return(updatedCandidate, nil).Once()
			},
			wantErr: false,
			resultCheck: func(t *testing.T, resp *CandidateResponse) {
				assert.Equal(t, "+9999999999", resp.Phone)
				assert.Equal(t, candidateID, resp.ID)
			},
		},

		"not found: candidate ID does not exist returns ErrCandidateNotFound": {
			req: &Request{
				ID: uuid.New(),
			},
			setupMock: func(m *MockCandidateRepo, req *Request) {
				m.On("GetByID", ctx, req.ID).
					Return(domain.Candidate{}, domain.ErrCandidateNotFound).
					Once()
			},
			wantErr: true,
			errCheck: func(err error) bool {
				return err == domain.ErrCandidateNotFound
			},
		},

		"repository error: database connection failed on update": {
			req: &Request{
				ID:    candidateID,
				Phone: stringPtr("+5555555555"),
			},
			setupMock: func(m *MockCandidateRepo, req *Request) {
				existingCandidate := domain.Candidate{
					ID:       candidateID,
					UserId:   userID,
					Phone:    "+1111111111",
					Telegram: "@user",
					City:     "City",
					Birthday: birthday,
				}
				m.On("GetByID", ctx, req.ID).Return(existingCandidate, nil).Once()
				m.On("Update", ctx, mock.AnythingOfType("domain.Candidate")).
					Return(domain.Candidate{}, domain.ErrCandidateNotFound).Once()
			},
			wantErr: true,
			errCheck: func(err error) bool {
				return err != nil
			},
		},

		"edge case: update all candidate fields simultaneously": {
			req: &Request{
				ID:       candidateID,
				UserID:   uuidPtr(uuid.New()),
				Phone:    stringPtr("+7777777777"),
				Telegram: stringPtr("@newuser"),
				City:     stringPtr("New City"),
				Birthday: timePtr(time.Date(1985, 5, 20, 0, 0, 0, 0, time.UTC)),
			},
			setupMock: func(m *MockCandidateRepo, req *Request) {
				existingCandidate := domain.Candidate{
					ID:       candidateID,
					UserId:   userID,
					Phone:    "+1111111111",
					Telegram: "@olduser",
					City:     "Old City",
					Birthday: birthday,
				}
				m.On("GetByID", ctx, req.ID).Return(existingCandidate, nil).Once()

				updatedCandidate := domain.Candidate{
					ID:       candidateID,
					UserId:   *req.UserID,
					Phone:    "+7777777777",
					Telegram: "@newuser",
					City:     "New City",
					Birthday: *req.Birthday,
				}
				m.On("Update", ctx, mock.AnythingOfType("domain.Candidate")).
					Return(updatedCandidate, nil).Once()
			},
			wantErr: false,
			resultCheck: func(t *testing.T, resp *CandidateResponse) {
				assert.Equal(t, "+7777777777", resp.Phone)
				assert.Equal(t, "@newuser", resp.Telegram)
				assert.Equal(t, "New City", resp.City)
			},
		},

		"edge case: update with empty strings to clear fields": {
			req: &Request{
				ID:       candidateID,
				Phone:    stringPtr(""),
				Telegram: stringPtr(""),
				City:     stringPtr(""),
			},
			setupMock: func(m *MockCandidateRepo, req *Request) {
				existingCandidate := domain.Candidate{
					ID:       candidateID,
					UserId:   userID,
					Phone:    "+1111111111",
					Telegram: "@user",
					City:     "City",
					Birthday: birthday,
				}
				m.On("GetByID", ctx, req.ID).Return(existingCandidate, nil).Once()

				clearedCandidate := domain.Candidate{
					ID:       candidateID,
					UserId:   userID,
					Phone:    "",
					Telegram: "",
					City:     "",
					Birthday: birthday,
				}
				m.On("Update", ctx, mock.MatchedBy(func(c domain.Candidate) bool {
					return c.Phone == "" && c.Telegram == "" && c.City == ""
				})).Return(clearedCandidate, nil).Once()
			},
			wantErr: false,
			resultCheck: func(t *testing.T, resp *CandidateResponse) {
				assert.Equal(t, "", resp.Phone)
				assert.Equal(t, "", resp.Telegram)
				assert.Equal(t, "", resp.City)
			},
		},

		"edge case: update request with no fields to update": {
			req: &Request{
				ID: candidateID,
				// No fields to update
			},
			setupMock: func(m *MockCandidateRepo, req *Request) {
				existingCandidate := domain.Candidate{
					ID:       candidateID,
					UserId:   userID,
					Phone:    "+1111111111",
					Telegram: "@user",
					City:     "City",
					Birthday: birthday,
				}
				m.On("GetByID", ctx, req.ID).Return(existingCandidate, nil).Once()
				m.On("Update", ctx, mock.AnythingOfType("domain.Candidate")).
					Return(existingCandidate, nil).Once()
			},
			wantErr: false,
			resultCheck: func(t *testing.T, resp *CandidateResponse) {
				assert.Equal(t, "+1111111111", resp.Phone)
				assert.Equal(t, "@user", resp.Telegram)
			},
		},

		"edge case: update with international values and special characters": {
			req: &Request{
				ID:       candidateID,
				Phone:    stringPtr("+33 (1) 42 68 53 00"),
				Telegram: stringPtr("@user-paris_2024"),
				City:     stringPtr("Paris, Île-de-France, France"),
			},
			setupMock: func(m *MockCandidateRepo, req *Request) {
				existingCandidate := domain.Candidate{
					ID:       candidateID,
					UserId:   userID,
					Phone:    "+1111111111",
					Telegram: "@user",
					City:     "City",
					Birthday: birthday,
				}
				m.On("GetByID", ctx, req.ID).Return(existingCandidate, nil).Once()

				updatedCandidate := domain.Candidate{
					ID:       candidateID,
					UserId:   userID,
					Phone:    "+33 (1) 42 68 53 00",
					Telegram: "@user-paris_2024",
					City:     "Paris, Île-de-France, France",
					Birthday: birthday,
				}
				m.On("Update", ctx, mock.AnythingOfType("domain.Candidate")).
					Return(updatedCandidate, nil).Once()
			},
			wantErr: false,
			resultCheck: func(t *testing.T, resp *CandidateResponse) {
				assert.Equal(t, "+33 (1) 42 68 53 00", resp.Phone)
			},
		},

		"edge case: update with future birth date": {
			req: &Request{
				ID:       candidateID,
				Birthday: timePtr(time.Now().AddDate(5, 0, 0)),
			},
			setupMock: func(m *MockCandidateRepo, req *Request) {
				existingCandidate := domain.Candidate{
					ID:       candidateID,
					UserId:   userID,
					Phone:    "+1111111111",
					Telegram: "@user",
					City:     "City",
					Birthday: birthday,
				}
				m.On("GetByID", ctx, req.ID).Return(existingCandidate, nil).Once()

				updatedCandidate := domain.Candidate{
					ID:       candidateID,
					UserId:   userID,
					Phone:    "+1111111111",
					Telegram: "@user",
					City:     "City",
					Birthday: *req.Birthday,
				}
				m.On("Update", ctx, mock.AnythingOfType("domain.Candidate")).
					Return(updatedCandidate, nil).Once()
			},
			wantErr: false,
			resultCheck: func(t *testing.T, resp *CandidateResponse) {
				assert.Equal(t, candidateID, resp.ID)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockCandidateRepo)
			tt.setupMock(mockRepo, tt.req)

			uc := New(mockRepo)

			// Execute
			result, err := uc.Execute(ctx, tt.req)

			// Verify
			if tt.wantErr {
				require.Error(t, err)
				require.True(t, tt.errCheck(err), "error check failed: %v", err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.resultCheck != nil {
					tt.resultCheck(t, result)
				}
			}

			// Cleanup: Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func uuidPtr(u uuid.UUID) *uuid.UUID {
	return &u
}

func timePtr(t time.Time) *time.Time {
	return &t
}
