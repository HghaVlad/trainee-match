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
	domain_errors "github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
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

type memRepoMock struct {
	mock.Mock
}

func (m *memRepoMock) Get(ctx context.Context, userID, companyID uuid.UUID) (*domain.CompanyMember, error) {
	res := m.Called(ctx, userID, companyID)

	if c := res.Get(0); c != nil {
		return c.(*domain.CompanyMember), res.Error(1)
	}

	return nil, res.Error(1)
}

func TestUsecase_Execute_HappyPath(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	compRepo := new(companyRepoMock)
	memRepo := new(memRepoMock)
	txManager := new(FakeTxManager)

	req := &create_vacancy.Request{
		CompanyID:   uuid.New(),
		Title:       "Go Dack Dev",
		Description: "Go Back dev pretty much",
		WorkFormat:  value_types.WorkFormatHybrid,
	}

	memRepo.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(&domain.CompanyMember{Role: value_types.CompanyRoleRecruiter}, nil).Once()

	vacRepo.On("Create", mock.Anything, mock.Anything).
		Return(nil).Once()

	compRepo.On("IncrementOpenVacancies", mock.Anything, mock.Anything).
		Return(nil).Once()

	uc := create_vacancy.NewUsecase(vacRepo, compRepo, memRepo, txManager)

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

	resp, err := uc.Execute(context.Background(), req, identity)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.True(t, txManager.Called)
	vacRepo.AssertExpectations(t)
	compRepo.AssertExpectations(t)
}

func TestUsecase_Execute_AuthErr(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	compRepo := new(companyRepoMock)
	txManager := new(FakeTxManager)

	req := &create_vacancy.Request{
		CompanyID:   uuid.New(),
		Title:       "Go Dack Dev",
		Description: "Go Back dev pretty much",
		WorkFormat:  value_types.WorkFormatHybrid,
	}

	tests := []struct {
		name     string
		identity uc_common.Identity
		memRepo  *memRepoMock
		resErr   error
	}{
		{
			name:     "Global Role Wrong",
			identity: uc_common.Identity{Role: uc_common.RoleCandidate},
			resErr:   domain_errors.ErrHrRoleRequired,
		},
		{
			name:     "HR is not member of company",
			identity: uc_common.Identity{Role: uc_common.RoleHR},
			memRepo: func() *memRepoMock {
				memRepo := new(memRepoMock)
				memRepo.On("Get", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, domain_errors.ErrCompanyMemberNotFound).Once()
				return memRepo
			}(),
			resErr: domain_errors.ErrCompanyMemberRequired,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uc := create_vacancy.NewUsecase(vacRepo, compRepo, test.memRepo, txManager)
			_, err := uc.Execute(context.Background(), req, test.identity)
			assert.Equal(t, test.resErr, err)
		})
	}
}

func TestUsecase_Execute_VacCreateFail(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	compRepo := new(companyRepoMock)
	memRepo := new(memRepoMock)
	txManager := new(FakeTxManager)

	req := &create_vacancy.Request{
		CompanyID:   uuid.New(),
		Title:       "Go Dack Dev",
		Description: "Go Back dev pretty much",
		WorkFormat:  value_types.WorkFormatHybrid,
	}

	memRepo.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(&domain.CompanyMember{Role: value_types.CompanyRoleRecruiter}, nil).Once()

	vacRepo.On("Create", mock.Anything, mock.Anything).
		Return(errors.New("some err")).Once()

	uc := create_vacancy.NewUsecase(vacRepo, compRepo, memRepo, txManager)

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

	_, err := uc.Execute(context.Background(), req, identity)

	require.Error(t, err)
	assert.True(t, txManager.Called)
	vacRepo.AssertExpectations(t)
	compRepo.AssertNotCalled(t, "IncrementOpenVacancies", mock.Anything, mock.Anything)
}

func TestUsecase_Execute_ValidateErr(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	compRepo := new(companyRepoMock)
	memRepo := new(memRepoMock)
	txManager := new(FakeTxManager)

	invalidReq := &create_vacancy.Request{
		CompanyID: uuid.New(),
	}

	uc := create_vacancy.NewUsecase(vacRepo, compRepo, memRepo, txManager)

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

	_, err := uc.Execute(context.Background(), invalidReq, identity)

	require.Error(t, err)
	vacRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	compRepo.AssertNotCalled(t, "IncrementOpenVacancies", mock.Anything, mock.Anything)
}
