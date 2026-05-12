package refresh_token

import (
	"context"
	"errors"
	"testing"

	"github.com/Nerzal/gocloak/v13"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/usecase/refresh_token/mocks"
)

func TestExecute(t *testing.T) {
	ctx := context.Background()
	validReq := &Request{RefreshToken: "refresh"}
	var (
		errAuth = errors.New("auth error")
	)
	validToken := &gocloak.JWT{AccessToken: "access", RefreshToken: "refresh"}

	tests := []struct {
		name          string
		req           *Request
		mockSetup     func(*mocks.MockAuthRepo)
		expectedToken *gocloak.JWT
		expectedError error
	}{
		{
			name: "valid request",
			req:  validReq,
			mockSetup: func(repo *mocks.MockAuthRepo) {
				repo.On("RefreshToken", mock.Anything, validReq.RefreshToken).Return(validToken, nil).Once()
			},
			expectedToken: validToken,
		},
		{
			name: "auth error",
			req:  validReq,
			mockSetup: func(repo *mocks.MockAuthRepo) {
				repo.On("RefreshToken", mock.Anything, validReq.RefreshToken).Return((*gocloak.JWT)(nil), errAuth).Once()
			},
			expectedToken: nil,
			expectedError: errAuth,
		},
		{
			name:          "context canceled",
			req:           validReq,
			expectedToken: nil,
			expectedError: context.Canceled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockAuthRepo(t)

			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}

			usecase := New(repo)
			testCtx := ctx
			if errors.Is(tt.expectedError, context.Canceled) {
				var cancel context.CancelFunc
				testCtx, cancel = context.WithCancel(ctx)
				cancel()
			}

			token, err := usecase.Execute(testCtx, tt.req)
			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectedError)
				require.Equal(t, tt.expectedToken, token)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedToken, token)
			}

			repo.AssertExpectations(t)
		})
	}
}

