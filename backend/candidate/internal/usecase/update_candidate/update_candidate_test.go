package update_candidate

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/update_candidate/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func stringPtr(s string) *string     { return &s }
func timePtr(t time.Time) *time.Time { return &t }

func TestExecute(t *testing.T) {
	ctx := context.Background()
	candidateID := uuid.New()
	userID := uuid.New()
	birthday := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)

	validRequest := &Request{
		Phone:    stringPtr("+1234567890"),
		Telegram: stringPtr("@valid_user"),
		City:     stringPtr("Valid City"),
		Birthday: timePtr(birthday),
	}

	var (
		errGetByUserID = errors.New("get by id db error")
		errUpdate      = errors.New("update db error")
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
				existing := domain.Candidate{ID: candidateID, UserId: userID, Phone: "+1111111111", Telegram: "@old", City: "Old City", Birthday: birthday}
				repo.On("GetByUserID", ctx, userID).Return(existing, nil).Once()
				updated := domain.Candidate{ID: candidateID, UserId: userID, Phone: "+1234567890", Telegram: "@valid_user", City: "Valid City", Birthday: birthday}
				repo.On("Update", ctx, mock.AnythingOfType("domain.Candidate")).Return(updated, nil).Once()
			},
			expectedID:    candidateID,
			expectedError: nil,
		},
		{
			name: "empty phone",
			request: &Request{
				Phone:    stringPtr(""),
				Telegram: stringPtr("@valid_user"),
				City:     stringPtr("Valid City"),
				Birthday: timePtr(birthday),
			},
			mockSetup: func(repo *mocks.CandidateRepo) {
				existing := domain.Candidate{ID: candidateID, UserId: userID, Phone: "+1111111111", Telegram: "@old", City: "Old City", Birthday: birthday}
				repo.On("GetByUserID", ctx, userID).Return(existing, nil).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrInvalidPhoneFormat,
		},
		{
			name: "invalid phone format",
			request: &Request{
				Phone:    stringPtr("invalid_phone"),
				Telegram: stringPtr("@valid_user"),
				City:     stringPtr("Valid City"),
				Birthday: timePtr(birthday),
			},
			mockSetup: func(repo *mocks.CandidateRepo) {
				existing := domain.Candidate{ID: candidateID, UserId: userID, Phone: "+1111111111", Telegram: "@old", City: "Old City", Birthday: birthday}
				repo.On("GetByUserID", ctx, userID).Return(existing, nil).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrInvalidPhoneFormat,
		},
		{
			name: "phone is not unique",
			request: &Request{
				Phone:    stringPtr("+1234567890"),
				Telegram: stringPtr("@valid_user"),
				City:     stringPtr("Valid City"),
				Birthday: timePtr(birthday),
			},
			mockSetup: func(repo *mocks.CandidateRepo) {
				existing := domain.Candidate{ID: candidateID, UserId: userID}
				repo.On("GetByUserID", ctx, userID).Return(existing, nil).Once()
				repo.On("Update", ctx, mock.AnythingOfType("domain.Candidate")).Return(domain.Candidate{}, domain.ErrPhoneAlreadyExists).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrPhoneAlreadyExists,
		},
		{
			name:    "phone is same as existing",
			request: &Request{Phone: stringPtr("+1111111111")},
			mockSetup: func(repo *mocks.CandidateRepo) {
				existing := domain.Candidate{ID: candidateID, UserId: userID, Phone: "+1111111111", Telegram: "@user", City: "City", Birthday: birthday}
				repo.On("GetByUserID", ctx, userID).Return(existing, nil).Once()
				updated := domain.Candidate{ID: candidateID, UserId: userID, Phone: "+1111111111", Telegram: "@user", City: "City", Birthday: birthday}
				repo.On("Update", ctx, mock.AnythingOfType("domain.Candidate")).Return(updated, nil).Once()
			},
			expectedID:    candidateID,
			expectedError: nil,
		},
		{
			name: "empty telegram",
			request: &Request{
				Phone:    stringPtr("+1234567890"),
				Telegram: stringPtr(""),
				City:     stringPtr("Valid City"),
				Birthday: timePtr(birthday),
			},
			mockSetup: func(repo *mocks.CandidateRepo) {
				existing := domain.Candidate{ID: candidateID, UserId: userID, Phone: "+1111111111", Telegram: "@old", City: "Old City", Birthday: birthday}
				repo.On("GetByUserID", ctx, userID).Return(existing, nil).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrInvalidTelegramFormat,
		},
		{
			name: "invalid telegram format",
			request: &Request{
				Phone:    stringPtr("+1234567890"),
				Telegram: stringPtr("invalid_telegram"),
				City:     stringPtr("Valid City"),
				Birthday: timePtr(birthday),
			},
			mockSetup: func(repo *mocks.CandidateRepo) {
				existing := domain.Candidate{ID: candidateID, UserId: userID, Phone: "+1111111111", Telegram: "@old", City: "Old City", Birthday: birthday}
				repo.On("GetByUserID", ctx, userID).Return(existing, nil).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrInvalidTelegramFormat,
		},
		{
			name: "telegram is not unique",
			request: &Request{
				Phone:    stringPtr("+1234567890"),
				Telegram: stringPtr("@valid_user"),
				City:     stringPtr("Valid City"),
				Birthday: timePtr(birthday),
			},
			mockSetup: func(repo *mocks.CandidateRepo) {
				existing := domain.Candidate{ID: candidateID, UserId: userID}
				repo.On("GetByUserID", ctx, userID).Return(existing, nil).Once()
				repo.On("Update", ctx, mock.AnythingOfType("domain.Candidate")).Return(domain.Candidate{}, domain.ErrTelegramAlreadyExists).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrTelegramAlreadyExists,
		},
		{
			name:    "telegram is same as existing",
			request: &Request{Telegram: stringPtr("@user")},
			mockSetup: func(repo *mocks.CandidateRepo) {
				existing := domain.Candidate{ID: candidateID, UserId: userID, Phone: "+1111111111", Telegram: "@user", City: "City", Birthday: birthday}
				repo.On("GetByUserID", ctx, userID).Return(existing, nil).Once()
				repo.On("Update", ctx, mock.AnythingOfType("domain.Candidate")).Return(existing, nil).Once()
			},
			expectedID:    candidateID,
			expectedError: nil,
		},
		{
			name: "city is empty",
			request: &Request{
				Phone:    stringPtr("+1234567890"),
				Telegram: stringPtr("@valid_user"),
				City:     stringPtr(""),
				Birthday: timePtr(birthday),
			},
			mockSetup: func(repo *mocks.CandidateRepo) {
				existing := domain.Candidate{ID: candidateID, UserId: userID, Phone: "+1111111111", Telegram: "@old", City: "Old City", Birthday: birthday}
				repo.On("GetByUserID", ctx, userID).Return(existing, nil).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrInvalidCityFormat,
		},
		{
			name: "birthday in the future",
			request: &Request{
				Phone:    stringPtr("+1234567890"),
				Telegram: stringPtr("@valid_user"),
				City:     stringPtr("Valid City"),
				Birthday: timePtr(time.Now().Add(24 * time.Hour)),
			},
			mockSetup: func(repo *mocks.CandidateRepo) {
				existing := domain.Candidate{ID: candidateID, UserId: userID, Phone: "+1111111111", Telegram: "@old", City: "Old City", Birthday: birthday}
				repo.On("GetByUserID", ctx, userID).Return(existing, nil).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: domain.ErrBirthdayInFuture,
		},
		{
			name:    "GetByUserID returns repo error",
			request: validRequest,
			mockSetup: func(repo *mocks.CandidateRepo) {
				repo.On("GetByUserID", ctx, userID).Return(domain.Candidate{}, errGetByUserID).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: errGetByUserID,
		},
		{
			name:    "Update returns repo error",
			request: &Request{Phone: stringPtr("+5555555555")},
			mockSetup: func(repo *mocks.CandidateRepo) {
				existing := domain.Candidate{ID: candidateID, UserId: userID, Phone: "+1111111111", Telegram: "@user", City: "City", Birthday: birthday}
				repo.On("GetByUserID", ctx, userID).Return(existing, nil).Once()
				repo.On("Update", ctx, mock.Anything).Return(domain.Candidate{}, errUpdate).Once()
			},
			expectedID:    uuid.Nil,
			expectedError: errUpdate,
		},
		{
			name:    "no fields",
			request: &Request{},
			mockSetup: func(repo *mocks.CandidateRepo) {
				existing := domain.Candidate{ID: candidateID, UserId: userID, Phone: "+1111111111", Telegram: "@user", City: "City", Birthday: birthday}
				repo.On("GetByUserID", ctx, userID).Return(existing, nil).Once()
				repo.On("Update", ctx, mock.Anything).Return(existing, nil).Once()
			},
			expectedID:    candidateID,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.CandidateRepo{}
			tt.mockSetup(repo)

			uc := New(repo)
			resp, err := uc.Execute(ctx, userID, tt.request)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.expectedError), "expected %v got %v", tt.expectedError, err)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tt.expectedID, resp.ID)
			}

			repo.AssertExpectations(t)
		})
	}
}
