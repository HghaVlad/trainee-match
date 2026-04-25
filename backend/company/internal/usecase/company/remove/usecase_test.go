package remove_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/remove"
)

type repoMock struct {
	mock.Mock
}

func (r *repoMock) Delete(ctx context.Context, id uuid.UUID) error {
	return r.Called(ctx, id).Error(0)
}

type cacheMock struct {
	mock.Mock
}

func (r *cacheMock) Del(ctx context.Context, id uuid.UUID) {
	r.Called(ctx, id)
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

func TestUsecase_Execute_OK(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)
	memRepo := new(memRepoMock)

	repo.On("Delete", mock.Anything, mock.Anything).
		Return(nil).Once()

	memRepo.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(&member.CompanyMember{Role: member.CompanyRoleAdmin}, nil).Once()

	cache.On("Del", mock.Anything, mock.Anything).Once()

	uc := remove.NewUsecase(repo, memRepo, cache)

	idenity := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	err := uc.Execute(context.Background(), uuid.New(), idenity)

	require.NoError(t, err)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestUsecase_Execute_CompanyRepoErr(t *testing.T) {
	repo := new(repoMock)
	cache := new(cacheMock)
	memRepo := new(memRepoMock)

	memRepo.On("Get", mock.Anything, mock.Anything, mock.Anything).
		Return(&member.CompanyMember{Role: member.CompanyRoleAdmin}, nil).Once()

	repo.On("Delete", mock.Anything, mock.Anything).
		Return(errors.New("some member err")).Once()

	uc := remove.NewUsecase(repo, memRepo, cache)

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	err := uc.Execute(context.Background(), uuid.New(), ident)

	require.Error(t, err)
	repo.AssertExpectations(t)
	cache.AssertNotCalled(t, "Del", mock.Anything, mock.Anything)
}
