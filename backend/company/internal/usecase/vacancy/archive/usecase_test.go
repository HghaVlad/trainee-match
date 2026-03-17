package archive_vacancy_test

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
	archive_vacancy "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/archive"
)

type vacancyRepoMock struct {
	mock.Mock
}

func (m *vacancyRepoMock) GetByID(ctx context.Context, vacID uuid.UUID, compID uuid.UUID) (*domain.Vacancy, error) {
	args := m.Called(ctx, vacID, compID)
	if vac := args.Get(0); vac != nil {
		return vac.(*domain.Vacancy), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *vacancyRepoMock) Archive(ctx context.Context, vacID uuid.UUID, compID uuid.UUID) error {
	return m.Called(ctx, vacID, compID).Error(0)
}

type companyRepoMock struct {
	mock.Mock
}

func (m *companyRepoMock) DecrementOpenVacancies(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

type memberRepoMock struct {
	mock.Mock
}

func (m *memberRepoMock) Get(ctx context.Context, userID, companyID uuid.UUID) (*domain.CompanyMember, error) {
	args := m.Called(ctx, userID, companyID)
	if member := args.Get(0); member != nil {
		return member.(*domain.CompanyMember), args.Error(1)
	}
	return nil, args.Error(1)
}

type cacheRepoMock struct {
	mock.Mock
}

func (m *cacheRepoMock) Del(ctx context.Context, id uuid.UUID) {
	m.Called(ctx, id)
}

type fakeTxManager struct {
	called bool
}

func (f *fakeTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	f.called = true
	return fn(ctx)
}

func TestUsecase_Execute_ArchivesPublishedVacancy(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	compRepo := new(companyRepoMock)
	memRepo := new(memberRepoMock)
	vacCache := new(cacheRepoMock)
	pubVacCache := new(cacheRepoMock)
	compCache := new(cacheRepoMock)
	txManager := new(fakeTxManager)

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}
	compID := uuid.New()
	vacID := uuid.New()

	memRepo.On("Get", mock.Anything, identity.UserID, compID).
		Return(&domain.CompanyMember{}, nil).Once()
	vacRepo.On("GetByID", mock.Anything, vacID, compID).
		Return(&domain.Vacancy{ID: vacID, CompanyID: compID, Status: value_types.VacancyStatusPublished}, nil).Once()
	vacRepo.On("Archive", mock.Anything, vacID, compID).
		Return(nil).Once()
	compRepo.On("DecrementOpenVacancies", mock.Anything, compID).
		Return(nil).Once()
	vacCache.On("Del", mock.Anything, vacID).Once()
	pubVacCache.On("Del", mock.Anything, vacID).Once()
	compCache.On("Del", mock.Anything, compID).Once()

	uc := archive_vacancy.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, pubVacCache, compCache)

	err := uc.Execute(context.Background(), compID, vacID, identity)

	require.NoError(t, err)
	assert.True(t, txManager.called)
	vacRepo.AssertExpectations(t)
	compRepo.AssertExpectations(t)
	memRepo.AssertExpectations(t)
	vacCache.AssertExpectations(t)
	pubVacCache.AssertExpectations(t)
	compCache.AssertExpectations(t)
}

func TestUsecase_Execute_ArchivesDraftWithoutCounterUpdate(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	compRepo := new(companyRepoMock)
	memRepo := new(memberRepoMock)
	vacCache := new(cacheRepoMock)
	pubVacCache := new(cacheRepoMock)
	compCache := new(cacheRepoMock)
	txManager := new(fakeTxManager)

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}
	compID := uuid.New()
	vacID := uuid.New()

	memRepo.On("Get", mock.Anything, identity.UserID, compID).
		Return(&domain.CompanyMember{}, nil).Once()
	vacRepo.On("GetByID", mock.Anything, vacID, compID).
		Return(&domain.Vacancy{ID: vacID, CompanyID: compID, Status: value_types.VacancyStatusDraft}, nil).Once()
	vacRepo.On("Archive", mock.Anything, vacID, compID).
		Return(nil).Once()
	vacCache.On("Del", mock.Anything, vacID).Once()
	pubVacCache.On("Del", mock.Anything, vacID).Once()
	compCache.On("Del", mock.Anything, compID).Once()

	uc := archive_vacancy.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, pubVacCache, compCache)

	err := uc.Execute(context.Background(), compID, vacID, identity)

	require.NoError(t, err)
	compRepo.AssertNotCalled(t, "DecrementOpenVacancies", mock.Anything, mock.Anything)
}

func TestUsecase_Execute_AlreadyArchived_NoOp(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	compRepo := new(companyRepoMock)
	memRepo := new(memberRepoMock)
	vacCache := new(cacheRepoMock)
	pubVacCache := new(cacheRepoMock)
	compCache := new(cacheRepoMock)
	txManager := new(fakeTxManager)

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}
	compID := uuid.New()
	vacID := uuid.New()

	memRepo.On("Get", mock.Anything, identity.UserID, compID).
		Return(&domain.CompanyMember{}, nil).Once()
	vacRepo.On("GetByID", mock.Anything, vacID, compID).
		Return(&domain.Vacancy{ID: vacID, CompanyID: compID, Status: value_types.VacancyStatusArchived}, nil).Once()

	uc := archive_vacancy.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, pubVacCache, compCache)

	err := uc.Execute(context.Background(), compID, vacID, identity)

	require.NoError(t, err)
	vacRepo.AssertNotCalled(t, "Archive", mock.Anything, mock.Anything, mock.Anything)
	compRepo.AssertNotCalled(t, "DecrementOpenVacancies", mock.Anything, mock.Anything)
	vacCache.AssertNotCalled(t, "Del", mock.Anything, mock.Anything)
	pubVacCache.AssertNotCalled(t, "Del", mock.Anything, mock.Anything)
	compCache.AssertNotCalled(t, "Del", mock.Anything, mock.Anything)
}

func TestUsecase_Execute_AuthErr(t *testing.T) {
	compID := uuid.New()
	vacID := uuid.New()

	t.Run("hr role required", func(t *testing.T) {
		vacRepo := new(vacancyRepoMock)
		compRepo := new(companyRepoMock)
		memRepo := new(memberRepoMock)
		vacCache := new(cacheRepoMock)
		pubVacCache := new(cacheRepoMock)
		compCache := new(cacheRepoMock)
		txManager := new(fakeTxManager)

		uc := archive_vacancy.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, pubVacCache, compCache)

		identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleCandidate}
		err := uc.Execute(context.Background(), compID, vacID, identity)

		assert.ErrorIs(t, err, domain_errors.ErrHrRoleRequired)
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
		memRepo.On("Get", mock.Anything, identity.UserID, compID).
			Return(nil, domain_errors.ErrCompanyMemberNotFound).Once()

		uc := archive_vacancy.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, pubVacCache, compCache)

		err := uc.Execute(context.Background(), compID, vacID, identity)

		assert.ErrorIs(t, err, domain_errors.ErrCompanyMemberRequired)
		vacRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestUsecase_Execute_RepoErr(t *testing.T) {
	compID := uuid.New()
	vacID := uuid.New()
	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

	t.Run("get vacancy", func(t *testing.T) {
		vacRepo := new(vacancyRepoMock)
		compRepo := new(companyRepoMock)
		memRepo := new(memberRepoMock)
		vacCache := new(cacheRepoMock)
		pubVacCache := new(cacheRepoMock)
		compCache := new(cacheRepoMock)
		txManager := new(fakeTxManager)

		memRepo.On("Get", mock.Anything, identity.UserID, compID).
			Return(&domain.CompanyMember{}, nil).Once()
		vacRepo.On("GetByID", mock.Anything, vacID, compID).
			Return(nil, errors.New("db err")).Once()

		uc := archive_vacancy.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, pubVacCache, compCache)

		err := uc.Execute(context.Background(), compID, vacID, identity)

		assert.EqualError(t, err, "db err")
	})

	t.Run("archive", func(t *testing.T) {
		vacRepo := new(vacancyRepoMock)
		compRepo := new(companyRepoMock)
		memRepo := new(memberRepoMock)
		vacCache := new(cacheRepoMock)
		pubVacCache := new(cacheRepoMock)
		compCache := new(cacheRepoMock)
		txManager := new(fakeTxManager)

		memRepo.On("Get", mock.Anything, identity.UserID, compID).
			Return(&domain.CompanyMember{}, nil).Once()
		vacRepo.On("GetByID", mock.Anything, vacID, compID).
			Return(&domain.Vacancy{Status: value_types.VacancyStatusPublished}, nil).Once()
		vacRepo.On("Archive", mock.Anything, vacID, compID).
			Return(errors.New("db err")).Once()

		uc := archive_vacancy.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, pubVacCache, compCache)

		err := uc.Execute(context.Background(), compID, vacID, identity)

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

		memRepo.On("Get", mock.Anything, identity.UserID, compID).
			Return(&domain.CompanyMember{}, nil).Once()
		vacRepo.On("GetByID", mock.Anything, vacID, compID).
			Return(&domain.Vacancy{Status: value_types.VacancyStatusPublished}, nil).Once()
		vacRepo.On("Archive", mock.Anything, vacID, compID).
			Return(nil).Once()
		compRepo.On("DecrementOpenVacancies", mock.Anything, compID).
			Return(errors.New("db err")).Once()

		uc := archive_vacancy.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, pubVacCache, compCache)

		err := uc.Execute(context.Background(), compID, vacID, identity)

		assert.EqualError(t, err, "db err")
		vacCache.AssertNotCalled(t, "Del", mock.Anything, mock.Anything)
		pubVacCache.AssertNotCalled(t, "Del", mock.Anything, mock.Anything)
		compCache.AssertNotCalled(t, "Del", mock.Anything, mock.Anything)
	})
}
