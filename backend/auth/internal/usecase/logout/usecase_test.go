package logout_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/usecase/logout"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/usecase/logout/mocks"
)

func TestExecute(t *testing.T) {
	ctx := context.Background()
	validReq := &logout.Request{Token: "refresh"}
	var (
		errAuth = errors.New("auth error")
	)

	tests := []struct {
		name          string
		req           *logout.Request
		mockSetup     func(*mocks.MockAuthRepo)
		expectedError error
	}{
		{
			name: "valid logout.Request",
			req:  validReq,
			mockSetup: func(repo *mocks.MockAuthRepo) {
				repo.On("Logout", mock.Anything, validReq.Token).Return(nil).Once()
			},
		},
		{
			name: "auth error",
			req:  validReq,
			mockSetup: func(repo *mocks.MockAuthRepo) {
				repo.On("Logout", mock.Anything, validReq.Token).Return(errAuth).Once()
			},
			expectedError: errAuth,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockAuthRepo(t)

			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			usecase := logout.New(repo)
			testCtx := ctx
			if errors.Is(tt.expectedError, context.Canceled) {
				var cancel context.CancelFunc
				testCtx, cancel = context.WithCancel(ctx)
				cancel()
			}

			err := usecase.Execute(testCtx, tt.req)
			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}
