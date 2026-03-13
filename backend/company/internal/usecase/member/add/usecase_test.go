package add_member_test

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
	add_member "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/add"
)

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

func (m *memberRepoMock) Create(ctx context.Context, member *domain.CompanyMember) error {
	return m.Called(ctx, member).Error(0)
}

func TestUsecase_ExecuteOK(t *testing.T) {
	repo := new(memberRepoMock)
	uc := add_member.NewUsecase(repo)

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}
	req := &add_member.Request{
		CompanyID: uuid.New(),
		UserID:    uuid.New(),
		Role:      value_types.CompanyRoleRecruiter,
	}

	repo.On("Get", mock.Anything, identity.UserID, req.CompanyID).
		Return(&domain.CompanyMember{Role: value_types.CompanyRoleAdmin}, nil).Once()
	repo.On("Create", mock.Anything, mock.MatchedBy(func(member *domain.CompanyMember) bool {
		return member.UserID == req.UserID &&
			member.CompanyID == req.CompanyID &&
			member.Role == req.Role
	})).Return(nil).Once()

	err := uc.Execute(context.Background(), req, identity)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUsecase_ExecuteAuthErr(t *testing.T) {
	repo := new(memberRepoMock)
	uc := add_member.NewUsecase(repo)

	req := &add_member.Request{
		CompanyID: uuid.New(),
		UserID:    uuid.New(),
		Role:      value_types.CompanyRoleRecruiter,
	}

	t.Run("global role required", func(t *testing.T) {
		identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleCandidate}

		err := uc.Execute(context.Background(), req, identity)

		assert.ErrorIs(t, err, domain_errors.ErrHrRoleRequired)
		repo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
		repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("company member required", func(t *testing.T) {
		identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

		repo.On("Get", mock.Anything, identity.UserID, req.CompanyID).
			Return(nil, domain_errors.ErrCompanyMemberNotFound).Once()

		err := uc.Execute(context.Background(), req, identity)

		assert.ErrorIs(t, err, domain_errors.ErrCompanyMemberRequired)
		repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("admin company role required", func(t *testing.T) {
		identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

		repo.On("Get", mock.Anything, identity.UserID, req.CompanyID).
			Return(&domain.CompanyMember{Role: value_types.CompanyRoleRecruiter}, nil).Once()

		err := uc.Execute(context.Background(), req, identity)

		assert.ErrorIs(t, err, domain_errors.ErrInsufficientRoleInCompany)
		repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})
}

func TestUsecase_ExecuteValidationAndRepoErr(t *testing.T) {
	t.Run("invalid user id", func(t *testing.T) {
		repo := new(memberRepoMock)
		uc := add_member.NewUsecase(repo)

		req := &add_member.Request{
			CompanyID: uuid.New(),
			UserID:    uuid.Nil,
			Role:      value_types.CompanyRoleRecruiter,
		}
		identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

		err := uc.Execute(context.Background(), req, identity)

		assert.ErrorIs(t, err, domain_errors.ErrInvalidUserID)
		repo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("invalid role", func(t *testing.T) {
		repo := new(memberRepoMock)
		uc := add_member.NewUsecase(repo)

		req := &add_member.Request{
			CompanyID: uuid.New(),
			UserID:    uuid.New(),
			Role:      "owner",
		}
		identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

		err := uc.Execute(context.Background(), req, identity)

		assert.ErrorIs(t, err, domain_errors.ErrInvalidCompanyMemberRole)
		repo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("repo err", func(t *testing.T) {
		repo := new(memberRepoMock)
		uc := add_member.NewUsecase(repo)

		identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}
		req := &add_member.Request{
			CompanyID: uuid.New(),
			UserID:    uuid.New(),
			Role:      value_types.CompanyRoleAdmin,
		}

		repo.On("Get", mock.Anything, identity.UserID, req.CompanyID).
			Return(&domain.CompanyMember{Role: value_types.CompanyRoleAdmin}, nil).Once()
		repo.On("Create", mock.Anything, mock.Anything).
			Return(errors.New("db err")).Once()

		err := uc.Execute(context.Background(), req, identity)

		assert.EqualError(t, err, "db err")
		repo.AssertExpectations(t)
	})
}
