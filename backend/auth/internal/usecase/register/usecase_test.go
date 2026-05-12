package register

import (
	"context"
	"errors"
	"testing"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/usecase/register/mocks"
)

func TestExecute(t *testing.T) {
	ctx := context.Background()
	validReq := &Request{
		FirstName: "Ivan",
		LastName:  "Petrov",
		Email:     "ivan.petrov@example.com",
		Username:  "ivanpetrov",
		Password:  "password123",
		Role:      "Candidate",
	}
	var (
		errDB     = errors.New("db error")
		errOutbox = errors.New("outbox error")
	)

	userID := uuid.New()

	tests := []struct {
		name          string
		req           *Request
		mockSetup     func(*mocks.MockAuthRepo, *mocks.MockOutboxWriter)
		expectedID    uuid.UUID
		expectedError error
	}{
		{
			name: "valid request",
			req:  validReq,
			mockSetup: func(repo *mocks.MockAuthRepo, writer *mocks.MockOutboxWriter) {
				repo.On("CreateUser", mock.Anything, mock.Anything, validReq.Password).Return(userID.String(), nil).Once()
				writer.On("WriteUserCreated", mock.Anything, mock.Anything).Return(nil).Once()
			},
			expectedID: userID,
		},
		{
			name: "invalid email",
			req: func() *Request {
				req := *validReq
				req.Email = "not-an-email"
				return &req
			}(),
			expectedID:    uuid.Nil,
			expectedError: domain.ErrInvalidEmail,
		},
		{
			name: "short password",
			req: func() *Request {
				req := *validReq
				req.Password = "short"
				return &req
			}(),
			expectedID:    uuid.Nil,
			expectedError: domain.ErrInvalidPassword,
		},
		{
			name: "missing name",
			req: func() *Request {
				req := *validReq
				req.FirstName = ""
				return &req
			}(),
			expectedID:    uuid.Nil,
			expectedError: domain.ErrInvalidName,
		},
		{
			name: "duplicate email error",
			req:  validReq,
			mockSetup: func(repo *mocks.MockAuthRepo, writer *mocks.MockOutboxWriter) {
				repo.On("CreateUser", mock.Anything, mock.Anything, validReq.Password).Return("", domain.ErrEmailAlreadyExists).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrEmailAlreadyExists,
		},
		{
			name: "db error",
			req:  validReq,
			mockSetup: func(repo *mocks.MockAuthRepo, writer *mocks.MockOutboxWriter) {
				repo.On("CreateUser", mock.Anything, mock.Anything, validReq.Password).Return("", errDB).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: errDB,
		},
		{
			name: "outbox error",
			req:  validReq,
			mockSetup: func(repo *mocks.MockAuthRepo, writer *mocks.MockOutboxWriter) {
				repo.On("CreateUser", mock.Anything, mock.Anything, validReq.Password).Return(userID.String(), nil).Once()
				writer.On("WriteUserCreated", mock.Anything, mock.Anything).Return(errOutbox).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: errOutbox,
		},
		{
			name:          "context canceled",
			req:           validReq,
			expectedID:    uuid.Nil,
			expectedError: context.Canceled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockAuthRepo(t)
			writer := mocks.NewMockOutboxWriter(t)

			if tt.mockSetup != nil {
				tt.mockSetup(repo, writer)
			}

			usecase := New(repo, writer)
			testCtx := ctx
			if errors.Is(tt.expectedError, context.Canceled) {
				var cancel context.CancelFunc
				testCtx, cancel = context.WithCancel(ctx)
				cancel()
			}

			id, err := usecase.Execute(testCtx, tt.req)
			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
				require.Equal(t, tt.expectedID, id)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedID, id)
			}

			repo.AssertExpectations(t)
			writer.AssertExpectations(t)
		})
	}
}

