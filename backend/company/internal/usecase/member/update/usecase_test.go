package update_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/update"
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

func (m *memberRepoMock) UpdateRole(ctx context.Context, userID, companyID uuid.UUID, role member.CompanyRole) error {
	return m.Called(ctx, userID, companyID, role).Error(0)
}

func TestUsecase_ExecuteOK(t *testing.T) {
	repo := new(memberRepoMock)
	uc := update.NewUsecase(repo)

	ident := identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}
	req := &update.Request{
		CompanyID: uuid.New(),
		UserID:    uuid.New(),
		Role:      member.CompanyRoleAdmin,
	}

	repo.On("Get", mock.Anything, ident.UserID, req.CompanyID).
		Return(&member.CompanyMember{Role: member.CompanyRoleAdmin}, nil).Once()
	repo.On("UpdateRole", mock.Anything, req.UserID, req.CompanyID, req.Role).
		Return(nil).Once()

	err := uc.Execute(context.Background(), req, ident)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUsecase_ExecuteAuthErr(t *testing.T) {
	repo := new(memberRepoMock)
	uc := update.NewUsecase(repo)

	req := &update.Request{
		CompanyID: uuid.New(),
		UserID:    uuid.New(),
		Role:      member.CompanyRoleRecruiter,
	}

	t.Run("global role required", func(t *testing.T) {
		ident := identity.Identity{UserID: uuid.New(), Role: identity.RoleCandidate}

		err := uc.Execute(context.Background(), req, ident)

		require.ErrorIs(t, err, identity.ErrHrRoleRequired)
		repo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
		repo.AssertNotCalled(t, "UpdateRole", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("company member required", func(t *testing.T) {
		ident := identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

		repo.On("Get", mock.Anything, ident.UserID, req.CompanyID).
			Return(nil, member.ErrCompanyMemberNotFound).Once()

		err := uc.Execute(context.Background(), req, ident)

		require.ErrorIs(t, err, member.ErrCompanyMemberRequired)
		repo.AssertNotCalled(t, "UpdateRole", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("admin company role required", func(t *testing.T) {
		ident := identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

		repo.On("Get", mock.Anything, ident.UserID, req.CompanyID).
			Return(&member.CompanyMember{Role: member.CompanyRoleRecruiter}, nil).Once()

		err := uc.Execute(context.Background(), req, ident)

		require.ErrorIs(t, err, member.ErrInsufficientRoleInCompany)
		repo.AssertNotCalled(t, "UpdateRole", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestUsecase_ExecuteValidationAndRepoErr(t *testing.T) {
	t.Run("invalid user id", func(t *testing.T) {
		repo := new(memberRepoMock)
		uc := update.NewUsecase(repo)

		req := &update.Request{
			CompanyID: uuid.New(),
			UserID:    uuid.Nil,
			Role:      member.CompanyRoleRecruiter,
		}
		ident := identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

		err := uc.Execute(context.Background(), req, ident)

		require.ErrorIs(t, err, member.ErrInvalidUserID)
		repo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("invalid role", func(t *testing.T) {
		repo := new(memberRepoMock)
		uc := update.NewUsecase(repo)

		req := &update.Request{
			CompanyID: uuid.New(),
			UserID:    uuid.New(),
			Role:      "owner",
		}
		ident := identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

		err := uc.Execute(context.Background(), req, ident)

		require.ErrorIs(t, err, member.ErrInvalidCompanyMemberRole)
		repo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("repo err", func(t *testing.T) {
		repo := new(memberRepoMock)
		uc := update.NewUsecase(repo)

		ident := identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}
		req := &update.Request{
			CompanyID: uuid.New(),
			UserID:    uuid.New(),
			Role:      member.CompanyRoleAdmin,
		}

		repo.On("Get", mock.Anything, ident.UserID, req.CompanyID).
			Return(&member.CompanyMember{Role: member.CompanyRoleAdmin}, nil).Once()
		repo.On("UpdateRole", mock.Anything, req.UserID, req.CompanyID, req.Role).
			Return(errors.New("db err")).Once()

		err := uc.Execute(context.Background(), req, ident)

		require.EqualError(t, err, "db err")
		repo.AssertExpectations(t)
	})
}
