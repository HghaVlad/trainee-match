package create_candidate

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/create_candidate/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
)

func TestExecute(t *testing.T) {
	ctx := context.Background()
	validRequest := &Request{
		UserID:   uuid.New(),
		Phone:    "+1234567890",
		Telegram: "@valid_user",
		City:     "Valid City",
		Birthday: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	validId := uuid.New()
	validUserID := uuid.New()
	var (
		errCreateDB = errors.New("create db error")
	)

	tests := []struct {
		name          string
		request       *Request
		mockSetup     func(repo *mocks.CandidateRepo)
		expectedID    uuid.UUID
		expectedError error
	}{
		{
			name:    "valid request",
			request: validRequest,
			mockSetup: func(repo *mocks.CandidateRepo) {
				repo.On("GetByUserID", ctx, validRequest.UserID).Return(domain.Candidate{}, domain.ErrCandidateNotFound).Maybe()
				repo.On("Create", ctx, mock.AnythingOfType("*domain.Candidate")).Return(validId, nil).Once()
			},
			expectedID:    validId,
			expectedError: nil,
		},
		{
			name: "empty phone",
			request: &Request{
				UserID:   validUserID,
				Phone:    "",
				Telegram: "@valid_user",
				City:     "Valid City",
				Birthday: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			mockSetup: func(repo *mocks.CandidateRepo) {
				repo.On("GetByUserID", ctx, validUserID).Return(domain.Candidate{}, domain.ErrCandidateNotFound).Maybe()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrInvalidPhoneFormat,
		},
		{
			name: "invalid phone format",
			request: &Request{
				UserID:   validUserID,
				Phone:    "invalid_phone",
				Telegram: "@valid_user",
				City:     "Valid City",
				Birthday: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			mockSetup: func(repo *mocks.CandidateRepo) {
				repo.On("GetByUserID", ctx, validUserID).Return(domain.Candidate{}, domain.ErrCandidateNotFound).Maybe()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrInvalidPhoneFormat,
		},
		{
			name: "phone is not unique",
			request: &Request{
				UserID:   validRequest.UserID,
				Phone:    "+1234567890",
				Telegram: "@valid_user",
				City:     "Valid City",
				Birthday: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			mockSetup: func(repo *mocks.CandidateRepo) {
				repo.On("GetByUserID", ctx, validRequest.UserID).Return(domain.Candidate{}, domain.ErrCandidateNotFound).Maybe()
				repo.On("Create", ctx, mock.AnythingOfType("*domain.Candidate")).Return(uuid.Nil, domain.ErrPhoneAlreadyExists).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrPhoneAlreadyExists,
		},
		{
			name: "empty telegram",
			request: &Request{
				UserID:   validUserID,
				Phone:    "+1234567890",
				Telegram: "",
				City:     "Valid City",
				Birthday: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			mockSetup: func(repo *mocks.CandidateRepo) {
				repo.On("GetByUserID", ctx, validUserID).Return(domain.Candidate{}, domain.ErrCandidateNotFound).Maybe()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrInvalidTelegramFormat,
		},
		{
			name: "invalid telegram format",
			request: &Request{
				UserID:   validUserID,
				Phone:    "+1234567890",
				Telegram: "invalid_telegram",
				City:     "Valid City",
				Birthday: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			mockSetup: func(repo *mocks.CandidateRepo) {
				repo.On("GetByUserID", ctx, validUserID).Return(domain.Candidate{}, domain.ErrCandidateNotFound).Maybe()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrInvalidTelegramFormat,
		},
		{
			name: "telegram is not unique",
			request: &Request{
				UserID:   validRequest.UserID,
				Phone:    "+1234567890",
				Telegram: "@valid_user",
				City:     "Valid City",
				Birthday: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			mockSetup: func(repo *mocks.CandidateRepo) {
				repo.On("GetByUserID", ctx, validRequest.UserID).Return(domain.Candidate{}, domain.ErrCandidateNotFound).Maybe()
				repo.On("Create", ctx, mock.AnythingOfType("*domain.Candidate")).Return(uuid.Nil, domain.ErrTelegramAlreadyExists).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrTelegramAlreadyExists,
		},
		{
			name: "city is empty",
			request: &Request{
				UserID:   validUserID,
				Phone:    "+1234567890",
				Telegram: "@valid_user",
				City:     "",
				Birthday: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			mockSetup: func(repo *mocks.CandidateRepo) {
				repo.On("GetByUserID", ctx, validUserID).Return(domain.Candidate{}, domain.ErrCandidateNotFound).Maybe()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrInvalidCityFormat,
		},
		{
			name: "birthday in the future",
			request: &Request{
				UserID:   validUserID,
				Phone:    "+1234567890",
				Telegram: "@valid_user",
				City:     "Valid City",
				Birthday: time.Now().Add(24 * time.Hour),
			},
			mockSetup: func(repo *mocks.CandidateRepo) {
				repo.On("GetByUserID", ctx, validUserID).Return(domain.Candidate{}, domain.ErrCandidateNotFound).Maybe()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrBirthdayInFuture,
		},
		{
			name:    "Create returns repo error",
			request: validRequest,
			mockSetup: func(repo *mocks.CandidateRepo) {
				repo.On("GetByUserID", ctx, validRequest.UserID).Return(domain.Candidate{}, domain.ErrCandidateNotFound).Maybe()
				repo.On("Create", ctx, mock.AnythingOfType("*domain.Candidate")).Return(uuid.Nil, errCreateDB).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: errCreateDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.CandidateRepo{}
			tt.mockSetup(repo)

			uc := New(repo)
			id, err := uc.Execute(ctx, tt.request)
			if tt.expectedError != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.expectedError), "expected error to be %v, got %v", tt.expectedError, err)
				require.Equal(t, uuid.Nil, id)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedID, id)
			}

			repo.AssertExpectations(t)
		})
	}

}
