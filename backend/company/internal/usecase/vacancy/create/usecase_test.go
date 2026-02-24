package create_vacancy_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/create"
)

type companyRepoMock struct {
	mock.Mock
}

func (m *companyRepoMock) IncrementOpenVacancies(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

type vacancyRepoMock struct {
	mock.Mock
}

func (m *vacancyRepoMock) Create(ctx context.Context, vacancy *domain.Vacancy) error {
	return m.Called(ctx, vacancy).Error(0)
}

type FakeTxManager struct {
	Called bool
}

func (f *FakeTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	f.Called = true
	return fn(ctx)
}

func TestUsecase_Execute_HappyPath(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	compRepo := new(companyRepoMock)
	txManager := new(FakeTxManager)

	req := &create_vacancy.Request{
		CompanyID:   uuid.New(),
		Title:       "Go Dack Dev",
		Description: "Go Back dev pretty much",
		WorkFormat:  value_types.WorkFormatHybrid,
	}

	vacRepo.On("Create", mock.Anything, mock.Anything).
		Return(nil).Once()

	compRepo.On("IncrementOpenVacancies", mock.Anything, mock.Anything).
		Return(nil).Once()

	uc := create_vacancy.NewUsecase(vacRepo, compRepo, txManager)

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.True(t, txManager.Called)
	vacRepo.AssertExpectations(t)
	compRepo.AssertExpectations(t)
}

func TestUsecase_Execute_VacCreateFail(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	compRepo := new(companyRepoMock)
	txManager := new(FakeTxManager)

	req := &create_vacancy.Request{
		CompanyID:   uuid.New(),
		Title:       "Go Dack Dev",
		Description: "Go Back dev pretty much",
		WorkFormat:  value_types.WorkFormatHybrid,
	}

	vacRepo.On("Create", mock.Anything, mock.Anything).
		Return(errors.New("some err")).Once()

	uc := create_vacancy.NewUsecase(vacRepo, compRepo, txManager)

	_, err := uc.Execute(context.Background(), req)

	require.Error(t, err)
	assert.True(t, txManager.Called)
	vacRepo.AssertExpectations(t)
	compRepo.AssertNotCalled(t, "IncrementOpenVacancies", mock.Anything, mock.Anything)
}

func TestUsecase_Execute_ValidateErr(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	compRepo := new(companyRepoMock)
	txManager := new(FakeTxManager)

	invalidReq := &create_vacancy.Request{
		CompanyID: uuid.New(),
	}

	uc := create_vacancy.NewUsecase(vacRepo, compRepo, txManager)

	_, err := uc.Execute(context.Background(), invalidReq)

	require.Error(t, err)
	vacRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	compRepo.AssertNotCalled(t, "IncrementOpenVacancies", mock.Anything, mock.Anything)
}
