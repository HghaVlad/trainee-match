package publish_test

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
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/publish"
)

type vacancyRepoMock struct {
	mock.Mock
}

func (m *vacancyRepoMock) GetByID(
	ctx context.Context,
	vacancyID uuid.UUID,
	companyID uuid.UUID,
) (*vacancy.Vacancy, error) {
	args := m.Called(ctx, vacancyID, companyID)
	if vac := args.Get(0); vac != nil {
		return vac.(*vacancy.Vacancy), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *vacancyRepoMock) Publish(ctx context.Context, vacID uuid.UUID, compID uuid.UUID) error {
	return m.Called(ctx, vacID, compID).Error(0)
}

type companyRepoMock struct {
	mock.Mock
}

func (m *companyRepoMock) IncrementOpenVacancies(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

type memberRepoMock struct {
	mock.Mock
}

func (m *memberRepoMock) Get(ctx context.Context, userID, companyID uuid.UUID) (*member.CompanyMember, error) {
	args := m.Called(ctx, userID, companyID)
	if memb := args.Get(0); memb != nil {
		return memb.(*member.CompanyMember), args.Error(1)
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

func TestUsecase_Execute_PublishesDraftVacancy(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	compRepo := new(companyRepoMock)
	memRepo := new(memberRepoMock)
	vacCache := new(cacheRepoMock)
	compCache := new(cacheRepoMock)
	txManager := new(fakeTxManager)

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}
	compID := uuid.New()
	vacID := uuid.New()

	memRepo.On("Get", mock.Anything, ident.UserID, compID).
		Return(&member.CompanyMember{}, nil).Once()
	vacRepo.On("GetByID", mock.Anything, vacID, compID).
		Return(&vacancy.Vacancy{ID: vacID, CompanyID: compID, Status: vacancy.StatusDraft}, nil).Once()
	vacRepo.On("Publish", mock.Anything, vacID, compID).
		Return(nil).Once()
	compRepo.On("IncrementOpenVacancies", mock.Anything, compID).
		Return(nil).Once()
	compCache.On("Del", mock.Anything, compID).Once()
	vacCache.On("Del", mock.Anything, vacID).Once()

	uc := publish.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, compCache)

	err := uc.Execute(context.Background(), compID, vacID, ident)

	require.NoError(t, err)
	assert.True(t, txManager.called)
	memRepo.AssertExpectations(t)
	vacRepo.AssertExpectations(t)
	compRepo.AssertExpectations(t)
	compCache.AssertExpectations(t)
	vacCache.AssertExpectations(t)
}

func TestUsecase_Execute_Alreadypublish_NoOp(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	compRepo := new(companyRepoMock)
	memRepo := new(memberRepoMock)
	vacCache := new(cacheRepoMock)
	compCache := new(cacheRepoMock)
	txManager := new(fakeTxManager)

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}
	compID := uuid.New()
	vacID := uuid.New()

	memRepo.On("Get", mock.Anything, ident.UserID, compID).
		Return(&member.CompanyMember{}, nil).Once()
	vacRepo.On("GetByID", mock.Anything, vacID, compID).
		Return(&vacancy.Vacancy{ID: vacID, CompanyID: compID, Status: vacancy.StatusPublished}, nil).Once()

	uc := publish.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, compCache)

	err := uc.Execute(context.Background(), compID, vacID, ident)

	require.NoError(t, err)
	assert.True(t, txManager.called)
	memRepo.AssertExpectations(t)
	vacRepo.AssertExpectations(t)
	vacRepo.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	compRepo.AssertNotCalled(t, "IncrementOpenVacancies", mock.Anything, mock.Anything)
	compCache.AssertNotCalled(t, "Del", mock.Anything, mock.Anything)
	vacCache.AssertNotCalled(t, "Del", mock.Anything, mock.Anything)
}

func TestUsecase_Execute_AuthErr(t *testing.T) {
	compID := uuid.New()
	vacID := uuid.New()

	t.Run("hr role required", func(t *testing.T) {
		vacRepo := new(vacancyRepoMock)
		compRepo := new(companyRepoMock)
		memRepo := new(memberRepoMock)
		vacCache := new(cacheRepoMock)
		compCache := new(cacheRepoMock)
		txManager := new(fakeTxManager)

		uc := publish.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, compCache)

		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleCandidate}
		err := uc.Execute(context.Background(), compID, vacID, ident)

		require.ErrorIs(t, err, identity.ErrHrRoleRequired)
		vacRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("company member required", func(t *testing.T) {
		vacRepo := new(vacancyRepoMock)
		compRepo := new(companyRepoMock)
		memRepo := new(memberRepoMock)
		vacCache := new(cacheRepoMock)
		compCache := new(cacheRepoMock)
		txManager := new(fakeTxManager)

		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}
		memRepo.On("Get", mock.Anything, ident.UserID, compID).
			Return(nil, member.ErrCompanyMemberNotFound).Once()

		uc := publish.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, compCache)

		err := uc.Execute(context.Background(), compID, vacID, ident)

		require.ErrorIs(t, err, member.ErrCompanyMemberRequired)
		memRepo.AssertExpectations(t)
		vacRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestUsecase_Execute_RepoErr(t *testing.T) {
	t.Run("get vacancy", func(t *testing.T) {
		vacRepo := new(vacancyRepoMock)
		compRepo := new(companyRepoMock)
		memRepo := new(memberRepoMock)
		vacCache := new(cacheRepoMock)
		compCache := new(cacheRepoMock)
		txManager := new(fakeTxManager)

		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}
		compID := uuid.New()
		vacID := uuid.New()

		memRepo.On("Get", mock.Anything, ident.UserID, compID).
			Return(&member.CompanyMember{}, nil).Once()
		vacRepo.On("GetByID", mock.Anything, vacID, compID).
			Return(nil, errors.New("db err")).Once()

		uc := publish.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, compCache)

		err := uc.Execute(context.Background(), compID, vacID, ident)

		require.EqualError(t, err, "db err")
		vacRepo.AssertExpectations(t)
		compRepo.AssertNotCalled(t, "IncrementOpenVacancies", mock.Anything, mock.Anything)
	})

	t.Run("publish", func(t *testing.T) {
		vacRepo := new(vacancyRepoMock)
		compRepo := new(companyRepoMock)
		memRepo := new(memberRepoMock)
		vacCache := new(cacheRepoMock)
		compCache := new(cacheRepoMock)
		txManager := new(fakeTxManager)

		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}
		compID := uuid.New()
		vacID := uuid.New()

		memRepo.On("Get", mock.Anything, ident.UserID, compID).
			Return(&member.CompanyMember{}, nil).Once()
		vacRepo.On("GetByID", mock.Anything, vacID, compID).
			Return(&vacancy.Vacancy{ID: vacID, CompanyID: compID, Status: vacancy.StatusDraft}, nil).Once()
		vacRepo.On("Publish", mock.Anything, vacID, compID).
			Return(errors.New("db err")).Once()

		uc := publish.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, compCache)

		err := uc.Execute(context.Background(), compID, vacID, ident)

		require.EqualError(t, err, "db err")
		vacRepo.AssertExpectations(t)
		compRepo.AssertNotCalled(t, "IncrementOpenVacancies", mock.Anything, mock.Anything)
	})

	t.Run("increment company vacancies", func(t *testing.T) {
		vacRepo := new(vacancyRepoMock)
		compRepo := new(companyRepoMock)
		memRepo := new(memberRepoMock)
		vacCache := new(cacheRepoMock)
		compCache := new(cacheRepoMock)
		txManager := new(fakeTxManager)

		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}
		compID := uuid.New()
		vacID := uuid.New()

		memRepo.On("Get", mock.Anything, ident.UserID, compID).
			Return(&member.CompanyMember{}, nil).Once()
		vacRepo.On("GetByID", mock.Anything, vacID, compID).
			Return(&vacancy.Vacancy{ID: vacID, CompanyID: compID, Status: vacancy.StatusDraft}, nil).Once()
		vacRepo.On("Publish", mock.Anything, vacID, compID).
			Return(nil).Once()
		compRepo.On("IncrementOpenVacancies", mock.Anything, compID).
			Return(errors.New("db err")).Once()

		uc := publish.NewUsecase(vacRepo, compRepo, memRepo, txManager, vacCache, compCache)

		err := uc.Execute(context.Background(), compID, vacID, ident)

		require.EqualError(t, err, "db err")
		compRepo.AssertExpectations(t)
		compCache.AssertNotCalled(t, "Del", mock.Anything, mock.Anything)
		vacCache.AssertNotCalled(t, "Del", mock.Anything, mock.Anything)
	})
}
