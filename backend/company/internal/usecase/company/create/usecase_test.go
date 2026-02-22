package create_company_test

import (
	"context"
	"testing"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/create"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type repoMock struct {
	mock.Mock
}

func (m *repoMock) Create(ctx context.Context, company *domain.Company) error {
	return m.Called(ctx, company).Error(0)
}

func TestUsecase_ExecuteOK(t *testing.T) {
	repo := new(repoMock)

	req := &create_company.Request{
		Name:        "Acme",
		Description: ptr("Hello!"),
	}

	repo.On("Create", mock.Anything, mock.Anything).
		Return(nil).Once()

	uc := create_company.NewUsecase(repo)

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	repo.AssertExpectations(t)
}

func TestUsecase_ExecuteValidateErr(t *testing.T) {
	repo := new(repoMock)
	uc := create_company.NewUsecase(repo)

	tests := []struct {
		name string
		req  create_company.Request
		err  error
	}{
		{
			name: "empty name",
			req: create_company.Request{
				Name: "",
			},
			err: domain_errors.ErrCompanyInvalidNameLen,
		},
		{
			name: "too long name",
			req: create_company.Request{
				Name: string(make([]byte, domain.MaxCompanyNameLen+1)),
			},
			err: domain_errors.ErrCompanyInvalidNameLen,
		},
		{
			name: "too long desc",
			req: create_company.Request{
				Name:        "Acme",
				Description: ptr(string(make([]byte, domain.MaxCompanyDescriptionLen+1))),
			},
			err: domain_errors.ErrCompanyInvalidDescriptionLen,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), &tt.req)

			assert.Equal(t, tt.err, err)
			repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
