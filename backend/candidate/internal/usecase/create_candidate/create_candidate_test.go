package create_candidate

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

func (m *MockCandidateRepo) Create(ctx context.Context, candidate *domain.Candidate) (uuid.UUID, error) {
	args := m.Called(ctx, candidate)
	var id uuid.UUID
	if args.Get(0) != nil {
		id = args.Get(0).(uuid.UUID)
	}
	return id, args.Error(1)
}

func TestExecute(t *testing.T) {
	ctx := context.Background()

	tests := map[string]struct {
		req         *Request
		setupMock   func(*MockCandidateRepo, *Request)
		wantErr     bool
		errCheck    func(error) bool
		resultCheck func(*testing.T, uuid.UUID)
	}{
		"happy path: valid candidate data creates successfully": {
			req: &Request{
				UserID:   uuid.New(),
				Phone:    "+1234567890",
				Telegram: "@testuser",
				City:     "New York",
				Birthday: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
			},
			setupMock: func(m *MockCandidateRepo, req *Request) {
				expectedID := uuid.New()
				// Mock expects Create to be called with any context and a candidate matching the request
				m.On("Create", mock.Anything, mock.MatchedBy(func(c *domain.Candidate) bool {
					return c.UserId == req.UserID &&
						c.Phone == req.Phone &&
						c.Telegram == req.Telegram &&
						c.City == req.City &&
						c.Birthday == req.Birthday
				})).Return(expectedID, nil).Once()
			},
			wantErr: false,
			resultCheck: func(t *testing.T, id uuid.UUID) {
				assert.NotEqual(t, uuid.Nil, id, "created candidate ID should not be nil")
			},
		},

		"repository error: database connection failure returns error": {
			req: &Request{
				UserID:   uuid.New(),
				Phone:    "+9876543210",
				Telegram: "@anotheruser",
				City:     "Boston",
				Birthday: time.Date(1985, 5, 20, 0, 0, 0, 0, time.UTC),
			},
			setupMock: func(m *MockCandidateRepo, req *Request) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.Candidate")).
					Return(uuid.Nil, domain.ErrCandidateNotFound).
					Once()
			},
			wantErr: true,
			errCheck: func(err error) bool {
				return err != nil
			},
		},

		"edge case: candidate with empty phone and telegram fields": {
			req: &Request{
				UserID:   uuid.New(),
				Phone:    "",
				Telegram: "",
				City:     "San Francisco",
				Birthday: time.Date(1988, 8, 10, 0, 0, 0, 0, time.UTC),
			},
			setupMock: func(m *MockCandidateRepo, req *Request) {
				expectedID := uuid.New()
				m.On("Create", mock.Anything, mock.MatchedBy(func(c *domain.Candidate) bool {
					return c.UserId == req.UserID &&
						c.Phone == "" &&
						c.Telegram == "" &&
						c.City == req.City
				})).Return(expectedID, nil).Once()
			},
			wantErr: false,
			resultCheck: func(t *testing.T, id uuid.UUID) {
				assert.NotEqual(t, uuid.Nil, id)
			},
		},

		"edge case: candidate with future birth date is accepted by usecase": {
			req: &Request{
				UserID:   uuid.New(),
				Phone:    "+1111111111",
				Telegram: "@future_user",
				City:     "Seattle",
				Birthday: time.Now().AddDate(1, 0, 0), // Future date
			},
			setupMock: func(m *MockCandidateRepo, req *Request) {
				expectedID := uuid.New()
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.Candidate")).
					Return(expectedID, nil).Once()
			},
			wantErr: false,
			resultCheck: func(t *testing.T, id uuid.UUID) {
				assert.NotEqual(t, uuid.Nil, id)
			},
		},

		"edge case: candidate with special characters in all fields": {
			req: &Request{
				UserID:   uuid.New(),
				Phone:    "+7 (999) 123-45-67 #123",
				Telegram: "@user_name-123-test",
				City:     "Санкт-Петербург, ул. Невского пр., 1",
				Birthday: time.Date(1992, 3, 15, 0, 0, 0, 0, time.UTC),
			},
			setupMock: func(m *MockCandidateRepo, req *Request) {
				expectedID := uuid.New()
				m.On("Create", mock.Anything, mock.MatchedBy(func(c *domain.Candidate) bool {
					return c.UserId == req.UserID &&
						c.Phone == req.Phone &&
						c.Telegram == req.Telegram &&
						c.City == req.City
				})).Return(expectedID, nil).Once()
			},
			wantErr: false,
			resultCheck: func(t *testing.T, id uuid.UUID) {
				assert.NotEqual(t, uuid.Nil, id)
			},
		},

		"edge case: nil UUID passed as UserID": {
			req: &Request{
				UserID:   uuid.Nil,
				Phone:    "+0000000000",
				Telegram: "@niluser",
				City:     "Unknown",
				Birthday: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			setupMock: func(m *MockCandidateRepo, req *Request) {
				expectedID := uuid.New()
				m.On("Create", mock.Anything, mock.MatchedBy(func(c *domain.Candidate) bool {
					return c.UserId == uuid.Nil
				})).Return(expectedID, nil).Once()
			},
			wantErr: false,
			resultCheck: func(t *testing.T, id uuid.UUID) {
				assert.NotEqual(t, uuid.Nil, id)
			},
		},

		"edge case: candidate with very long field values": {
			req: &Request{
				UserID:   uuid.New(),
				Phone:    "+1234567890" + "0123456789" + "0123456789",
				Telegram: "@" + string(make([]byte, 100)),
				City:     string(make([]byte, 200)),
				Birthday: time.Date(1975, 6, 20, 0, 0, 0, 0, time.UTC),
			},
			setupMock: func(m *MockCandidateRepo, req *Request) {
				expectedID := uuid.New()
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.Candidate")).
					Return(expectedID, nil).Once()
			},
			wantErr: false,
			resultCheck: func(t *testing.T, id uuid.UUID) {
				assert.NotEqual(t, uuid.Nil, id)
			},
		},

		"repository error: constraint violation returns error": {
			req: &Request{
				UserID:   uuid.New(),
				Phone:    "+2222222222",
				Telegram: "@constraint_test",
				City:     "Miami",
				Birthday: time.Date(1995, 7, 25, 0, 0, 0, 0, time.UTC),
			},
			setupMock: func(m *MockCandidateRepo, req *Request) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.Candidate")).
					Return(uuid.Nil, domain.ErrCandidateNotFound). // Simulating any repository error
					Once()
			},
			wantErr: true,
			errCheck: func(err error) bool {
				return err != nil
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
			gotID, err := uc.Execute(ctx, tt.req)

			// Verify
			if tt.wantErr {
				require.Error(t, err)
				require.True(t, tt.errCheck(err), "error check failed: %v", err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, gotID)
				if tt.resultCheck != nil {
					tt.resultCheck(t, gotID)
				}
			}

			// Cleanup: Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}
