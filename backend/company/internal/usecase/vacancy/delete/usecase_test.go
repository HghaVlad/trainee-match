package delete_vacancy_test

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
	uc_common "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	delete_vacancy "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/delete"
)

type vacancyRepoMock struct {
	mock.Mock
}

func (m *vacancyRepoMock) GetByID(ctx context.Context, vacancyID uuid.UUID, companyID uuid.UUID) (*domain.Vacancy, error) {
	args := m.Called(ctx, vacancyID, companyID)
	if vac := args.Get(0); vac != nil {
		return vac.(*domain.Vacancy), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *vacancyRepoMock) Delete(ctx context.Context, vacancyID uuid.UUID, companyID uuid.UUID) error {
	return m.Called(ctx, vacancyID, companyID).Error(0)
}

type companyRepoMock struct {
	mock.Mock
}

func (m *companyRepoMock) DecrementOpenVacancies(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

type cacheRepoMock struct {
	mock.Mock
}

func (m *cacheRepoMock) Del(ctx context.Context, id uuid.UUID) {
	m.Called(ctx, id)
}

type memberRepoMock struct {
	mock.Mock
}

func (m *memberRepoMock) Get(ctx context.Context, userID, companyID uuid.UUID) (*domain.CompanyMember, error) {
	res := m.Called(ctx, userID, companyID)
	if c := res.Get(0); c != nil {
		return c.(*domain.CompanyMember), res.Error(1)
	}
	return nil, res.Error(1)
}

type fakeTxManager struct {
	called bool
}

func (f *fakeTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	f.called = true
	return fn(ctx)
}

func TestUsecase_Execute_DeletesPublishedVacancy(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	compRepo := new(companyRepoMock)
	memRepo := new(memberRepoMock)
	vacCache := new(cacheRepoMock)
	pubVacCache := new(cacheRepoMock)
	compCache := new(cacheRepoMock)
	txManager := new(fakeTxManager)

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}
	companyID := uuid.New()
	vacancyID := uuid.New()

	memRepo.On("Get", mock.Anything, identity.UserID, companyID).
		Return(&domain.CompanyMember{}, nil).Once()
	vacRepo.On("GetByID", mock.Anything, vacancyID, companyID).
		Return(&domain.Vacancy{ID: vacancyID, CompanyID: companyID, Status: value_types.VacancyStatusPublished}, nil).Once()
	vacRepo.On("Delete", mock.Anything, vacancyID, companyID).
		Return(nil).Once()
	compRepo.On("DecrementOpenVacancies", mock.Anything, companyID).
		Return(nil).Once()
	vacCache.On("Del", mock.Anything, vacancyID).Once()
	pubVacCache.On("Del", mock.Anything, vacancyID).Once()
	compCache.On("Del", mock.Anything, companyID).Once()

	uc := delete_vacancy.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, pubVacCache, compCache)

	err := uc.Execute(context.Background(), vacancyID, companyID, identity)

	require.NoError(t, err)
	assert.True(t, txManager.called)
	memRepo.AssertExpectations(t)
	vacRepo.AssertExpectations(t)
	compRepo.AssertExpectations(t)
	vacCache.AssertExpectations(t)
	pubVacCache.AssertExpectations(t)
	compCache.AssertExpectations(t)
}

func TestUsecase_Execute_DeletesDraftWithoutCounterUpdate(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	compRepo := new(companyRepoMock)
	memRepo := new(memberRepoMock)
	vacCache := new(cacheRepoMock)
	pubVacCache := new(cacheRepoMock)
	compCache := new(cacheRepoMock)
	txManager := new(fakeTxManager)

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}
	companyID := uuid.New()
	vacancyID := uuid.New()

	memRepo.On("Get", mock.Anything, identity.UserID, companyID).
		Return(&domain.CompanyMember{}, nil).Once()
	vacRepo.On("GetByID", mock.Anything, vacancyID, companyID).
		Return(&domain.Vacancy{ID: vacancyID, CompanyID: companyID, Status: value_types.VacancyStatusDraft}, nil).Once()
	vacRepo.On("Delete", mock.Anything, vacancyID, companyID).
		Return(nil).Once()
	vacCache.On("Del", mock.Anything, vacancyID).Once()
	pubVacCache.On("Del", mock.Anything, vacancyID).Once()
	compCache.On("Del", mock.Anything, companyID).Once()

	uc := delete_vacancy.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, pubVacCache, compCache)

	err := uc.Execute(context.Background(), vacancyID, companyID, identity)

	require.NoError(t, err)
	compRepo.AssertNotCalled(t, "DecrementOpenVacancies", mock.Anything, mock.Anything)
}

func TestUsecase_Execute_AuthErr(t *testing.T) {
	companyID := uuid.New()
	vacancyID := uuid.New()

	t.Run("hr role required", func(t *testing.T) {
		vacRepo := new(vacancyRepoMock)
		compRepo := new(companyRepoMock)
		memRepo := new(memberRepoMock)
		vacCache := new(cacheRepoMock)
		pubVacCache := new(cacheRepoMock)
		compCache := new(cacheRepoMock)
		txManager := new(fakeTxManager)

		uc := delete_vacancy.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, pubVacCache, compCache)

		identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleCandidate}
		err := uc.Execute(context.Background(), vacancyID, companyID, identity)

		assert.ErrorIs(t, err, domain_errors.ErrHrRoleRequired)
		memRepo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
		vacRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("company member required", func(t *testing.T) {
		vacRepo := new(vacancyRepoMock)
		compRepo := new(companyRepoMock)
		memRepo := new(memberRepoMock)
		vacCache := new(cacheRepoMock)
		pubVacCache := new(cacheRepoMock)
		compCache := new(cacheRepoMock)
		txManager := new(fakeTxManager)

		identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}
		memRepo.On("Get", mock.Anything, identity.UserID, companyID).
			Return(nil, domain_errors.ErrCompanyMemberNotFound).Once()

		uc := delete_vacancy.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, pubVacCache, compCache)

		err := uc.Execute(context.Background(), vacancyID, companyID, identity)

		assert.ErrorIs(t, err, domain_errors.ErrCompanyMemberRequired)
		vacRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestUsecase_Execute_RepoErr(t *testing.T) {
	companyID := uuid.New()
	vacancyID := uuid.New()
	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

	t.Run("get vacancy", func(t *testing.T) {
		vacRepo := new(vacancyRepoMock)
		compRepo := new(companyRepoMock)
		memRepo := new(memberRepoMock)
		vacCache := new(cacheRepoMock)
		pubVacCache := new(cacheRepoMock)
		compCache := new(cacheRepoMock)
		txManager := new(fakeTxManager)

		memRepo.On("Get", mock.Anything, identity.UserID, companyID).
			Return(&domain.CompanyMember{}, nil).Once()
		vacRepo.On("GetByID", mock.Anything, vacancyID, companyID).
			Return(nil, errors.New("db err")).Once()

		uc := delete_vacancy.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, pubVacCache, compCache)

		err := uc.Execute(context.Background(), vacancyID, companyID, identity)

		assert.EqualError(t, err, "db err")
		vacRepo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("delete", func(t *testing.T) {
		vacRepo := new(vacancyRepoMock)
		compRepo := new(companyRepoMock)
		memRepo := new(memberRepoMock)
		vacCache := new(cacheRepoMock)
		pubVacCache := new(cacheRepoMock)
		compCache := new(cacheRepoMock)
		txManager := new(fakeTxManager)

		memRepo.On("Get", mock.Anything, identity.UserID, companyID).
			Return(&domain.CompanyMember{}, nil).Once()
		vacRepo.On("GetByID", mock.Anything, vacancyID, companyID).
			Return(&domain.Vacancy{Status: value_types.VacancyStatusPublished}, nil).Once()
		vacRepo.On("Delete", mock.Anything, vacancyID, companyID).
			Return(errors.New("db err")).Once()

		uc := delete_vacancy.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, pubVacCache, compCache)

		err := uc.Execute(context.Background(), vacancyID, companyID, identity)

		assert.EqualError(t, err, "db err")
		compRepo.AssertNotCalled(t, "DecrementOpenVacancies", mock.Anything, mock.Anything)
	})

	t.Run("decrement company vacancies", func(t *testing.T) {
		vacRepo := new(vacancyRepoMock)
		compRepo := new(companyRepoMock)
		memRepo := new(memberRepoMock)
		vacCache := new(cacheRepoMock)
		pubVacCache := new(cacheRepoMock)
		compCache := new(cacheRepoMock)
		txManager := new(fakeTxManager)

		memRepo.On("Get", mock.Anything, identity.UserID, companyID).
			Return(&domain.CompanyMember{}, nil).Once()
		vacRepo.On("GetByID", mock.Anything, vacancyID, companyID).
			Return(&domain.Vacancy{Status: value_types.VacancyStatusPublished}, nil).Once()
		vacRepo.On("Delete", mock.Anything, vacancyID, companyID).
			Return(nil).Once()
		compRepo.On("DecrementOpenVacancies", mock.Anything, companyID).
			Return(errors.New("db err")).Once()

		uc := delete_vacancy.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, pubVacCache, compCache)

		err := uc.Execute(context.Background(), vacancyID, companyID, identity)

		assert.EqualError(t, err, "db err")
		vacCache.AssertNotCalled(t, "Del", mock.Anything, mock.Anything)
		pubVacCache.AssertNotCalled(t, "Del", mock.Anything, mock.Anything)
		compCache.AssertNotCalled(t, "Del", mock.Anything, mock.Anything)
	})
}
