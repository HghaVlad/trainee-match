package delete_member_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domain "github.com/HghaVlad/trainee-match/backend/company/internal/domain/entities"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/delete"
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

func (m *memberRepoMock) Delete(ctx context.Context, userID, companyID uuid.UUID) error {
	return m.Called(ctx, userID, companyID).Error(0)
}

func TestUsecase_ExecuteOK(t *testing.T) {
	repo := new(memberRepoMock)
	uc := delete_member.NewUsecase(repo)

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}
	companyID := uuid.New()
	userID := uuid.New()

	repo.On("Get", mock.Anything, identity.UserID, companyID).
		Return(&domain.CompanyMember{Role: value_types.CompanyRoleAdmin}, nil).Once()
	repo.On("Delete", mock.Anything, userID, companyID).
		Return(nil).Once()

	err := uc.Execute(context.Background(), companyID, userID, identity)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUsecase_ExecuteAuthErr(t *testing.T) {
	repo := new(memberRepoMock)
	uc := delete_member.NewUsecase(repo)

	companyID := uuid.New()
	userID := uuid.New()

	t.Run("global role required", func(t *testing.T) {
		identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleCandidate}

		err := uc.Execute(context.Background(), companyID, userID, identity)

		assert.ErrorIs(t, err, domain_errors.ErrHrRoleRequired)
		repo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
		repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("company member required", func(t *testing.T) {
		identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

		repo.On("Get", mock.Anything, identity.UserID, companyID).
			Return(nil, domain_errors.ErrCompanyMemberNotFound).Once()

		err := uc.Execute(context.Background(), companyID, userID, identity)

		assert.ErrorIs(t, err, domain_errors.ErrCompanyMemberRequired)
		repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("admin company role required", func(t *testing.T) {
		identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

		repo.On("Get", mock.Anything, identity.UserID, companyID).
			Return(&domain.CompanyMember{Role: value_types.CompanyRoleRecruiter}, nil).Once()

		err := uc.Execute(context.Background(), companyID, userID, identity)

		assert.ErrorIs(t, err, domain_errors.ErrInsufficientRoleInCompany)
		repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestUsecase_ExecuteRepoErr(t *testing.T) {
	repo := new(memberRepoMock)
	uc := delete_member.NewUsecase(repo)

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}
	companyID := uuid.New()
	userID := uuid.New()

	repo.On("Get", mock.Anything, identity.UserID, companyID).
		Return(&domain.CompanyMember{Role: value_types.CompanyRoleAdmin}, nil).Once()
	repo.On("Delete", mock.Anything, userID, companyID).
		Return(errors.New("db err")).Once()

	err := uc.Execute(context.Background(), companyID, userID, identity)

	assert.EqualError(t, err, "db err")
	repo.AssertExpectations(t)
}
