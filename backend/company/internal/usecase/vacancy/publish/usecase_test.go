package publish_vacancy_test

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
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/publish"
)

type vacancyRepoMock struct {
	mock.Mock
}

func (m *vacancyRepoMock) Publish(ctx context.Context, compID uuid.UUID, vacID uuid.UUID) error {
	return m.Called(ctx, compID, vacID).Error(0)
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

func TestUsecase_ExecuteOK(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	memRepo := new(memberRepoMock)
	uc := publish_vacancy.NewUsecase(vacRepo, memRepo)

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}
	compID := uuid.New()
	vacID := uuid.New()

	memRepo.On("Get", mock.Anything, identity.UserID, compID).
		Return(&domain.CompanyMember{}, nil).Once()
	vacRepo.On("Publish", mock.Anything, compID, vacID).
		Return(nil).Once()

	err := uc.Execute(context.Background(), compID, vacID, identity)

	require.NoError(t, err)
	memRepo.AssertExpectations(t)
	vacRepo.AssertExpectations(t)
}

func TestUsecase_ExecuteAuthErr(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	memRepo := new(memberRepoMock)
	uc := publish_vacancy.NewUsecase(vacRepo, memRepo)

	compID := uuid.New()
	vacID := uuid.New()

	t.Run("hr role required", func(t *testing.T) {
		identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleCandidate}

		err := uc.Execute(context.Background(), compID, vacID, identity)

		assert.ErrorIs(t, err, domain_errors.ErrHrRoleRequired)
		memRepo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
		vacRepo.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("company member required", func(t *testing.T) {
		identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

		memRepo.On("Get", mock.Anything, identity.UserID, compID).
			Return(nil, domain_errors.ErrCompanyMemberNotFound).Once()

		err := uc.Execute(context.Background(), compID, vacID, identity)

		assert.ErrorIs(t, err, domain_errors.ErrCompanyMemberRequired)
		vacRepo.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
		memRepo.AssertExpectations(t)
	})
}

func TestUsecase_ExecuteRepoErr(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	memRepo := new(memberRepoMock)
	uc := publish_vacancy.NewUsecase(vacRepo, memRepo)

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}
	compID := uuid.New()
	vacID := uuid.New()

	memRepo.On("Get", mock.Anything, identity.UserID, compID).
		Return(&domain.CompanyMember{}, nil).Once()
	vacRepo.On("Publish", mock.Anything, compID, vacID).
		Return(errors.New("db err")).Once()

	err := uc.Execute(context.Background(), compID, vacID, identity)

	assert.EqualError(t, err, "db err")
	memRepo.AssertExpectations(t)
	vacRepo.AssertExpectations(t)
}
