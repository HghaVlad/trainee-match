package get_vacancy_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	domain_errors "github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	uc_common "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/get_by_id"
)

type repoMock struct {
	mock.Mock
}

func (m *repoMock) GetByID(ctx context.Context, id uuid.UUID, companyID uuid.UUID) (*domain.Vacancy, error) {
	args := m.Called(ctx, id, companyID)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Vacancy), args.Error(1)
	}
	return nil, args.Error(1)
}

type cacheMock struct {
	mock.Mock
}

func (m *cacheMock) Put(ctx context.Context, key uuid.UUID, val *domain.Vacancy, exp time.Duration) {
	m.Called(ctx, key, val, exp)
}

func (m *cacheMock) Get(ctx context.Context, id uuid.UUID) *domain.Vacancy {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Vacancy)
	}
	return nil
}

type memberRepoMock struct {
	mock.Mock
}

func (m *memberRepoMock) Get(ctx context.Context, userID, companyID uuid.UUID) (*domain.CompanyMember, error) {
	args := m.Called(ctx, userID, companyID)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.CompanyMember), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestUsecase_Execute_CacheHit(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)
	memberRepo := new(memberRepoMock)

	id := uuid.New()
	compID := uuid.New()
	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

	memberRepo.On("Get", mock.Anything, identity.UserID, compID).
		Return(&domain.CompanyMember{}, nil).Once()

	cache.On("Get", mock.Anything, mock.Anything).
		Return(&domain.Vacancy{ID: id, CompanyID: compID, Title: "Title"}).Once()

	uc := get_vacancy.NewUsecase(repo, cache, memberRepo)

	resp, err := uc.Execute(context.Background(), id, compID, identity)

	require.NoError(t, err)
	assert.Equal(t, resp.ID, id)
	assert.Equal(t, resp.Title, "Title")
	memberRepo.AssertExpectations(t)
	cache.AssertExpectations(t)
	repo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
}

func TestUsecase_Execute_CacheMiss(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)
	memberRepo := new(memberRepoMock)

	id := uuid.New()
	compID := uuid.New()
	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

	memberRepo.On("Get", mock.Anything, identity.UserID, compID).
		Return(&domain.CompanyMember{}, nil).Once()

	cache.On("Get", mock.Anything, mock.Anything).
		Return(nil).Once()

	repo.On("GetByID", mock.Anything, mock.Anything, mock.Anything).
		Return(&domain.Vacancy{ID: id, CompanyID: compID, Title: "Title"}, nil).Once()

	cache.On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once()

	uc := get_vacancy.NewUsecase(repo, cache, memberRepo)

	resp, err := uc.Execute(context.Background(), id, compID, identity)

	require.NoError(t, err)
	assert.Equal(t, resp.ID, id)
	assert.Equal(t, resp.Title, "Title")
	memberRepo.AssertExpectations(t)
	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestUsecase_Execute_RepoErr(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)
	memberRepo := new(memberRepoMock)

	id := uuid.New()
	compID := uuid.New()
	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

	memberRepo.On("Get", mock.Anything, identity.UserID, compID).
		Return(&domain.CompanyMember{}, nil).Once()

	cache.On("Get", mock.Anything, mock.Anything).
		Return(nil).Once()

	repo.On("GetByID", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("err: i. e. not found ")).Once()

	uc := get_vacancy.NewUsecase(repo, cache, memberRepo)

	_, err := uc.Execute(context.Background(), id, compID, identity)

	require.Error(t, err)
	memberRepo.AssertExpectations(t)
	cache.AssertExpectations(t)
	repo.AssertExpectations(t)
	cache.AssertNotCalled(t, "Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestUsecase_Execute_AuthErr(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)
	memberRepo := new(memberRepoMock)
	uc := get_vacancy.NewUsecase(repo, cache, memberRepo)

	id := uuid.New()
	compID := uuid.New()

	t.Run("hr role required", func(t *testing.T) {
		identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleCandidate}

		_, err := uc.Execute(context.Background(), id, compID, identity)

		assert.ErrorIs(t, err, domain_errors.ErrHrRoleRequired)
		memberRepo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
		cache.AssertNotCalled(t, "Get", mock.Anything, mock.Anything)
		repo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("company member required", func(t *testing.T) {
		identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

		memberRepo.On("Get", mock.Anything, identity.UserID, compID).
			Return(nil, domain_errors.ErrCompanyMemberNotFound).Once()

		_, err := uc.Execute(context.Background(), id, compID, identity)

		assert.ErrorIs(t, err, domain_errors.ErrCompanyMemberRequired)
		memberRepo.AssertExpectations(t)
		cache.AssertNotCalled(t, "Get", mock.Anything, mock.Anything)
		repo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything, mock.Anything)
	})
}
