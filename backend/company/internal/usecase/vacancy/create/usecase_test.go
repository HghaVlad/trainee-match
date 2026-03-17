package create_vacancy_test

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
	create_vacancy "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/create"
)

type vacancyRepoMock struct {
	mock.Mock
}

func (m *vacancyRepoMock) Create(ctx context.Context, vacancy *domain.Vacancy) error {
	return m.Called(ctx, vacancy).Error(0)
}

type memRepoMock struct {
	mock.Mock
}

func (m *memRepoMock) Get(ctx context.Context, userID, companyID uuid.UUID) (*domain.CompanyMember, error) {
	res := m.Called(ctx, userID, companyID)

	if c := res.Get(0); c != nil {
		return c.(*domain.CompanyMember), res.Error(1)
	}

	return nil, res.Error(1)
}

func TestUsecase_Execute_HappyPath(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	memRepo := new(memRepoMock)

	req := &create_vacancy.Request{
		CompanyID:   uuid.New(),
		Title:       "Go Backend Dev",
		Description: "Go backend dev pretty much",
		WorkFormat:  value_types.WorkFormatHybrid,
	}

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

	memRepo.On("Get", mock.Anything, identity.UserID, req.CompanyID).
		Return(&domain.CompanyMember{Role: value_types.CompanyRoleRecruiter}, nil).Once()

	vacRepo.On("Create", mock.Anything, mock.MatchedBy(func(v *domain.Vacancy) bool {
		return v.CompanyID == req.CompanyID &&
			v.CreatedBy == identity.UserID &&
			v.Title == req.Title &&
			v.Description == req.Description &&
			v.WorkFormat == req.WorkFormat &&
			v.Status == value_types.VacancyStatusDraft &&
			v.EmploymentType == value_types.EmploymentTypeInternship &&
			v.ID != uuid.Nil
	})).Return(nil).Once()

	uc := create_vacancy.NewUsecase(vacRepo, memRepo)

	resp, err := uc.Execute(context.Background(), req, identity)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEqual(t, uuid.Nil, resp.ID)
	memRepo.AssertExpectations(t)
	vacRepo.AssertExpectations(t)
}

func TestUsecase_Execute_UsesProvidedEmploymentType(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	memRepo := new(memRepoMock)

	employmentType := value_types.EmploymentTypePartTime
	req := &create_vacancy.Request{
		CompanyID:         uuid.New(),
		Title:             "Go Backend Dev",
		Description:       "Go backend dev pretty much",
		WorkFormat:        value_types.WorkFormatRemote,
		EmploymentType:    &employmentType,
		InternshipToOffer: true,
	}

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

	memRepo.On("Get", mock.Anything, identity.UserID, req.CompanyID).
		Return(&domain.CompanyMember{Role: value_types.CompanyRoleRecruiter}, nil).Once()

	vacRepo.On("Create", mock.Anything, mock.MatchedBy(func(v *domain.Vacancy) bool {
		return v.EmploymentType == value_types.EmploymentTypePartTime &&
			v.Status == value_types.VacancyStatusDraft &&
			v.InternshipToOffer
	})).Return(nil).Once()

	uc := create_vacancy.NewUsecase(vacRepo, memRepo)

	_, err := uc.Execute(context.Background(), req, identity)

	require.NoError(t, err)
	memRepo.AssertExpectations(t)
	vacRepo.AssertExpectations(t)
}

func TestUsecase_Execute_AuthErr(t *testing.T) {
	vacRepo := new(vacancyRepoMock)

	req := &create_vacancy.Request{
		CompanyID:   uuid.New(),
		Title:       "Go Backend Dev",
		Description: "Go backend dev pretty much",
		WorkFormat:  value_types.WorkFormatHybrid,
	}

	t.Run("global role wrong", func(t *testing.T) {
		memRepo := new(memRepoMock)
		uc := create_vacancy.NewUsecase(vacRepo, memRepo)

		_, err := uc.Execute(context.Background(), req, uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleCandidate})

		assert.ErrorIs(t, err, domain_errors.ErrHrRoleRequired)
		memRepo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
		vacRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("hr is not member of company", func(t *testing.T) {
		memRepo := new(memRepoMock)
		identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

		memRepo.On("Get", mock.Anything, identity.UserID, req.CompanyID).
			Return(nil, domain_errors.ErrCompanyMemberNotFound).Once()

		uc := create_vacancy.NewUsecase(vacRepo, memRepo)

		_, err := uc.Execute(context.Background(), req, identity)

		assert.ErrorIs(t, err, domain_errors.ErrCompanyMemberRequired)
		memRepo.AssertExpectations(t)
		vacRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})
}

func TestUsecase_Execute_VacCreateFail(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	memRepo := new(memRepoMock)

	req := &create_vacancy.Request{
		CompanyID:   uuid.New(),
		Title:       "Go Backend Dev",
		Description: "Go backend dev pretty much",
		WorkFormat:  value_types.WorkFormatHybrid,
	}

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

	memRepo.On("Get", mock.Anything, identity.UserID, req.CompanyID).
		Return(&domain.CompanyMember{Role: value_types.CompanyRoleRecruiter}, nil).Once()

	vacRepo.On("Create", mock.Anything, mock.Anything).
		Return(errors.New("some err")).Once()

	uc := create_vacancy.NewUsecase(vacRepo, memRepo)

	_, err := uc.Execute(context.Background(), req, identity)

	require.EqualError(t, err, "some err")
	memRepo.AssertExpectations(t)
	vacRepo.AssertExpectations(t)
}

func TestUsecase_Execute_ValidateErr(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	memRepo := new(memRepoMock)

	invalidReq := &create_vacancy.Request{
		CompanyID: uuid.New(),
	}

	uc := create_vacancy.NewUsecase(vacRepo, memRepo)

	identity := uc_common.Identity{UserID: uuid.New(), Role: uc_common.RoleHR}

	_, err := uc.Execute(context.Background(), invalidReq, identity)

	require.Error(t, err)
	memRepo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
	vacRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}
