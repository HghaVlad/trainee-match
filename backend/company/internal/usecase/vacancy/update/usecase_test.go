package update_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/update"
)

type repoMock struct {
	mock.Mock
}

func (m *repoMock) GetByID(ctx context.Context, vacancyID uuid.UUID, companyID uuid.UUID) (*vacancy.Vacancy, error) {
	args := m.Called(ctx, vacancyID, companyID)

	if v := args.Get(0); v != nil {
		return v.(*vacancy.Vacancy), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *repoMock) Update(ctx context.Context, v *vacancy.Vacancy) error {
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

func (m *memRepoMock) Get(ctx context.Context, userID, companyID uuid.UUID) (*member.CompanyMember, error) {
	res := m.Called(ctx, userID, companyID)

	if c := res.Get(0); c != nil {
		return c.(*member.CompanyMember), res.Error(1)
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

	req := &update.Request{
		VacancyID: vID,
		CompanyID: cID,
		Title:     ptr("New Title"),
	}

	vac := vacancy.Vacancy{
		ID:             vID,
		CompanyID:      cID,
		Title:          "Go dev",
		Description:    "Go back dev",
		WorkFormat:     vacancy.WorkFormatHybrid,
		EmploymentType: vacancy.EmploymentTypeInternship,
	}

	memRepo.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(&member.CompanyMember{Role: member.CompanyRoleAdmin}, nil).Once()

	repo.On("GetByID", mock.Anything, vID, cID).
		Return(&vac, nil).Once()

	repo.On("Update", mock.Anything, mock.Anything).
		Return(nil).Once()

	cache.On("Del", mock.Anything, mock.Anything).Once()

	uc := update.NewUsecase(repo, memRepo, cache, txManager)

	ident := identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	err := uc.Execute(context.Background(), req, ident)

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

	req := &update.Request{
		VacancyID: vID,
		CompanyID: cID,
		Title:     ptr("New Title"),
	}

	memRepo.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(&member.CompanyMember{Role: member.CompanyRoleAdmin}, nil).Once()

	repo.On("GetByID", mock.Anything, vID, cID).
		Return(nil, errors.New("repo get err")).Once()

	uc := update.NewUsecase(repo, memRepo, cache, txManager)

	ident := identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	err := uc.Execute(context.Background(), req, ident)

	require.Error(t, err)

	assert.True(t, txManager.called)
	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	cache.AssertNotCalled(t, "Del", mock.Anything, mock.Anything)
}

func ptr[T any](v T) *T {
	return &v
}
