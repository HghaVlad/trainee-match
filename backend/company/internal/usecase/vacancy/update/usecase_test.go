package update_vacancy_test

import (
	"context"
	"errors"
	"testing"

	uc_common "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/update"
)

type repoMock struct {
	mock.Mock
}

func (m *repoMock) GetByID(ctx context.Context, vacancyID uuid.UUID, companyID uuid.UUID) (*domain.Vacancy, error) {
	args := m.Called(ctx, vacancyID, companyID)

	if v := args.Get(0); v != nil {
		return v.(*domain.Vacancy), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *repoMock) Update(ctx context.Context, v *domain.Vacancy) error {
	return m.Called(ctx, v).Error(0)
}

type cacheMock struct {
	mock.Mock
}

func (m *cacheMock) Del(ctx context.Context, id uuid.UUID) {
	m.Called(ctx, id)
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

type FakeTxManager struct {
	called bool
}

func (m *FakeTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	m.called = true
	return fn(ctx)
}

func TestUsecase_Execute_HappyPath(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)
	txManager := new(FakeTxManager)
	memRepo := new(memRepoMock)

	vID := uuid.New()
	cID := uuid.New()

	req := &update_vacancy.Request{
		VacancyID: vID,
		CompanyID: cID,
		Title:     ptr("New Title"),
	}

	vac := domain.Vacancy{
		ID:             vID,
		CompanyID:      cID,
		Title:          "Go dev",
		Description:    "Go back dev",
		WorkFormat:     value_types.WorkFormatHybrid,
		EmploymentType: value_types.EmploymentTypeInternship,
	}

	memRepo.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(&domain.CompanyMember{Role: value_types.CompanyRoleAdmin}, nil).Once()

	repo.On("GetByID", mock.Anything, vID, cID).
		Return(&vac, nil).Once()

	repo.On("Update", mock.Anything, mock.Anything).
		Return(nil).Once()

	cache.On("Del", mock.Anything, mock.Anything).Once()

	uc := update_vacancy.NewUsecase(repo, memRepo, cache, txManager)

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

	err := uc.Execute(context.Background(), req, identity)

	require.NoError(t, err)
	assert.Equal(t, vac.ID, req.VacancyID)
	assert.Equal(t, vac.Title, *req.Title)

	assert.True(t, txManager.called)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestUsecase_Execute_GetErr(t *testing.T) {
	repo := new(repoMock)
	memRepo := new(memRepoMock)
	cache := new(cacheMock)
	txManager := new(FakeTxManager)

	vID := uuid.New()
	cID := uuid.New()

	req := &update_vacancy.Request{
		VacancyID: vID,
		CompanyID: cID,
		Title:     ptr("New Title"),
	}

	memRepo.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(&domain.CompanyMember{Role: value_types.CompanyRoleAdmin}, nil).Once()

	repo.On("GetByID", mock.Anything, vID, cID).
		Return(nil, errors.New("repo get err")).Once()

	uc := update_vacancy.NewUsecase(repo, memRepo, cache, txManager)

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

	err := uc.Execute(context.Background(), req, identity)

	require.Error(t, err)

	assert.True(t, txManager.called)
	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	cache.AssertNotCalled(t, "Del", mock.Anything, mock.Anything)
}

func ptr[T any](v T) *T {
	return &v
}
