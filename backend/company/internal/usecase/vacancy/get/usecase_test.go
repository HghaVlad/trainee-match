package get_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/get"
)

type repoMock struct {
	mock.Mock
}

func (m *repoMock) GetByID(ctx context.Context, id uuid.UUID, companyID uuid.UUID) (*vacancy.Vacancy, error) {
	args := m.Called(ctx, id, companyID)
	if args.Get(0) != nil {
		return args.Get(0).(*vacancy.Vacancy), args.Error(1)
	}
	return nil, args.Error(1)
}

type cacheMock struct {
	mock.Mock
}

func (m *cacheMock) Put(ctx context.Context, key uuid.UUID, val *vacancy.Vacancy, exp time.Duration) {
	m.Called(ctx, key, val, exp)
}

func (m *cacheMock) Get(ctx context.Context, id uuid.UUID) *vacancy.Vacancy {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*vacancy.Vacancy)
	}
	return nil
}

type memberRepoMock struct {
	mock.Mock
}

func (m *memberRepoMock) Get(ctx context.Context, userID, companyID uuid.UUID) (*member.CompanyMember, error) {
	args := m.Called(ctx, userID, companyID)
	if args.Get(0) != nil {
		return args.Get(0).(*member.CompanyMember), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestUsecase_Execute_CacheHit(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)
	memberRepo := new(memberRepoMock)

	id := uuid.New()
	compID := uuid.New()
	ident := identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	memberRepo.On("Get", mock.Anything, ident.UserID, compID).
		Return(&member.CompanyMember{}, nil).Once()

	cache.On("Get", mock.Anything, mock.Anything).
		Return(&vacancy.Vacancy{ID: id, CompanyID: compID, Title: "Title"}).Once()

	uc := get.NewUsecase(repo, cache, memberRepo)

	resp, err := uc.Execute(context.Background(), id, compID, ident)

	require.NoError(t, err)
	assert.Equal(t, resp.ID, id)
	assert.Equal(t, "Title", resp.Title)
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
	ident := identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	memberRepo.On("Get", mock.Anything, ident.UserID, compID).
		Return(&member.CompanyMember{}, nil).Once()

	cache.On("Get", mock.Anything, mock.Anything).
		Return(nil).Once()

	repo.On("GetByID", mock.Anything, mock.Anything, mock.Anything).
		Return(&vacancy.Vacancy{ID: id, CompanyID: compID, Title: "Title"}, nil).Once()

	cache.On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once()

	uc := get.NewUsecase(repo, cache, memberRepo)

	resp, err := uc.Execute(context.Background(), id, compID, ident)

	require.NoError(t, err)
	assert.Equal(t, resp.ID, id)
	assert.Equal(t, "Title", resp.Title)
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
	ident := identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	memberRepo.On("Get", mock.Anything, ident.UserID, compID).
		Return(&member.CompanyMember{}, nil).Once()

	cache.On("Get", mock.Anything, mock.Anything).
		Return(nil).Once()

	repo.On("GetByID", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("err: i. e. not found ")).Once()

	uc := get.NewUsecase(repo, cache, memberRepo)

	_, err := uc.Execute(context.Background(), id, compID, ident)

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
	uc := get.NewUsecase(repo, cache, memberRepo)

	id := uuid.New()
	compID := uuid.New()

	t.Run("hr role required", func(t *testing.T) {
		ident := identity.Identity{UserID: uuid.New(), Role: identity.RoleCandidate}

		_, err := uc.Execute(context.Background(), id, compID, ident)

		require.ErrorIs(t, err, identity.ErrHrRoleRequired)
		memberRepo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
		cache.AssertNotCalled(t, "Get", mock.Anything, mock.Anything)
		repo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("company member required", func(t *testing.T) {
		ident := identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

		memberRepo.On("Get", mock.Anything, ident.UserID, compID).
			Return(nil, member.ErrCompanyMemberNotFound).Once()

		_, err := uc.Execute(context.Background(), id, compID, ident)

		require.ErrorIs(t, err, member.ErrCompanyMemberRequired)
		memberRepo.AssertExpectations(t)
		cache.AssertNotCalled(t, "Get", mock.Anything, mock.Anything)
		repo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything, mock.Anything)
	})
}
