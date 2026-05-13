package getuser_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	getuser "github.com/HghaVlad/trainee-match/backend/auth/internal/usecase/get_user_me"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/domain"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/usecase/get_user_me/mocks"
)

func TestExecute(t *testing.T) {
	ctx := context.Background()
	validReq := &getuser.Request{Token: "access"}
	var (
		errAuth = errors.New("auth error")
		errRole = errors.New("role error")
	)

	baseUser := func() *domain.User {
		return &domain.User{
			ID:        "user-id",
			FirstName: "Ivan",
			LastName:  "Petrov",
			Email:     "ivan.petrov@example.com",
			Username:  "ivanpetrov",
		}
	}
	baseUserWithRole := func(role string) *domain.User {
		user := baseUser()
		user.Role = role
		return user
	}

	tests := []struct {
		name          string
		req           *getuser.Request
		mockSetup     func(*mocks.MockAuthRepo)
		expectedUser  *domain.User
		expectedError error
	}{
		{
			name: "valid getuser.Request",
			req:  validReq,
			mockSetup: func(repo *mocks.MockAuthRepo) {
				user := baseUser()
				repo.On("GetUserInfo", mock.Anything, validReq.Token).Return(user, nil).Once()
				repo.On("GetUserRole", mock.Anything, validReq.Token, user.ID).Return("Candidate", nil).Once()
			},
			expectedUser: baseUserWithRole("Candidate"),
		},
		{
			name: "get user error",
			req:  validReq,
			mockSetup: func(repo *mocks.MockAuthRepo) {
				repo.On("GetUserInfo", mock.Anything, validReq.Token).Return((*domain.User)(nil), errAuth).Once()
			},
			expectedUser:  nil,
			expectedError: errAuth,
		},
		{
			name: "get role error",
			req:  validReq,
			mockSetup: func(repo *mocks.MockAuthRepo) {
				user := baseUser()
				repo.On("GetUserInfo", mock.Anything, validReq.Token).Return(user, nil).Once()
				repo.On("GetUserRole", mock.Anything, validReq.Token, user.ID).Return("", errRole).Once()
			},
			expectedUser:  nil,
			expectedError: errRole,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockAuthRepo(t)

			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			usecase := getuser.New(repo)
			testCtx := ctx
			if errors.Is(tt.expectedError, context.Canceled) {
				var cancel context.CancelFunc
				testCtx, cancel = context.WithCancel(ctx)
				cancel()
			}

			result, err := usecase.Execute(testCtx, tt.req)
			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
				require.Equal(t, tt.expectedUser, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedUser, result)
			}

			repo.AssertExpectations(t)
		})
	}
}
