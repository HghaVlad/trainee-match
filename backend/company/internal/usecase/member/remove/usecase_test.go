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
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/remove"
)

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

func (m *memberRepoMock) Delete(ctx context.Context, userID, companyID uuid.UUID) error {
	return m.Called(ctx, userID, companyID).Error(0)
}

func TestUsecase_ExecuteOK(t *testing.T) {
	repo := new(memberRepoMock)
	uc := remove.NewUsecase(repo)

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}
	companyID := uuid.New()
	userID := uuid.New()

	repo.On("Get", mock.Anything, ident.UserID, companyID).
		Return(&member.CompanyMember{Role: member.CompanyRoleAdmin}, nil).Once()
	repo.On("Delete", mock.Anything, userID, companyID).
		Return(nil).Once()

	err := uc.Execute(context.Background(), companyID, userID, ident)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUsecase_ExecuteAuthErr(t *testing.T) {
	repo := new(memberRepoMock)
	uc := remove.NewUsecase(repo)

	companyID := uuid.New()
	userID := uuid.New()

	t.Run("global role required", func(t *testing.T) {
		iden := &identity.Identity{UserID: uuid.New(), Role: identity.RoleCandidate}

		err := uc.Execute(context.Background(), companyID, userID, iden)

		require.ErrorIs(t, err, identity.ErrHrRoleRequired)
		repo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
		repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("company member required", func(t *testing.T) {
		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

		repo.On("Get", mock.Anything, ident.UserID, companyID).
			Return(nil, member.ErrCompanyMemberNotFound).Once()

		err := uc.Execute(context.Background(), companyID, userID, ident)

		require.ErrorIs(t, err, member.ErrCompanyMemberRequired)
		repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("admin company role required", func(t *testing.T) {
		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

		repo.On("Get", mock.Anything, ident.UserID, companyID).
			Return(&member.CompanyMember{Role: member.CompanyRoleRecruiter}, nil).Once()

		err := uc.Execute(context.Background(), companyID, userID, ident)

		require.ErrorIs(t, err, member.ErrInsufficientRoleInCompany)
		repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestUsecase_ExecuteRepoErr(t *testing.T) {
	repo := new(memberRepoMock)
	uc := remove.NewUsecase(repo)

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}
	companyID := uuid.New()
	userID := uuid.New()

	repo.On("Get", mock.Anything, ident.UserID, companyID).
		Return(&member.CompanyMember{Role: member.CompanyRoleAdmin}, nil).Once()
	repo.On("Delete", mock.Anything, userID, companyID).
		Return(errors.New("db err")).Once()

	err := uc.Execute(context.Background(), companyID, userID, ident)

	require.EqualError(t, err, "db err")
	repo.AssertExpectations(t)
}
