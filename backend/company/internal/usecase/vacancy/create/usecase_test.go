package create_test

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
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/create"
)

type vacancyRepoMock struct {
	mock.Mock
}

func (m *vacancyRepoMock) Create(ctx context.Context, vacancy *vacancy.Vacancy) error {
	return m.Called(ctx, vacancy).Error(0)
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

func TestUsecase_Execute_HappyPath(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	memRepo := new(memRepoMock)

	req := &create.Request{
		CompanyID:   uuid.New(),
		Title:       "Go Backend Dev",
		Description: "Go backend dev pretty much",
		WorkFormat:  vacancy.WorkFormatHybrid,
	}

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	memRepo.On("Get", mock.Anything, ident.UserID, req.CompanyID).
		Return(&member.CompanyMember{Role: member.CompanyRoleRecruiter}, nil).Once()

	vacRepo.On("Create", mock.Anything, mock.MatchedBy(func(v *vacancy.Vacancy) bool {
		return v.CompanyID == req.CompanyID &&
			v.CreatedBy == ident.UserID &&
			v.Title == req.Title &&
			v.Description == req.Description &&
			v.WorkFormat == req.WorkFormat &&
			v.Status == vacancy.StatusDraft &&
			v.EmploymentType == vacancy.EmploymentTypeInternship &&
			v.ID != uuid.Nil
	})).Return(nil).Once()

	uc := create.NewUsecase(vacRepo, memRepo)

	resp, err := uc.Execute(context.Background(), req, ident)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEqual(t, uuid.Nil, resp.ID)
	memRepo.AssertExpectations(t)
	vacRepo.AssertExpectations(t)
}

func TestUsecase_Execute_UsesProvidedEmploymentType(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	memRepo := new(memRepoMock)

	employmentType := vacancy.EmploymentTypePartTime
	req := &create.Request{
		CompanyID:         uuid.New(),
		Title:             "Go Backend Dev",
		Description:       "Go backend dev pretty much",
		WorkFormat:        vacancy.WorkFormatRemote,
		EmploymentType:    &employmentType,
		InternshipToOffer: true,
	}

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	memRepo.On("Get", mock.Anything, ident.UserID, req.CompanyID).
		Return(&member.CompanyMember{Role: member.CompanyRoleRecruiter}, nil).Once()

	vacRepo.On("Create", mock.Anything, mock.MatchedBy(func(v *vacancy.Vacancy) bool {
		return v.EmploymentType == vacancy.EmploymentTypePartTime &&
			v.Status == vacancy.StatusDraft &&
			v.InternshipToOffer
	})).Return(nil).Once()

	uc := create.NewUsecase(vacRepo, memRepo)

	_, err := uc.Execute(context.Background(), req, ident)

	require.NoError(t, err)
	memRepo.AssertExpectations(t)
	vacRepo.AssertExpectations(t)
}

func TestUsecase_Execute_AuthErr(t *testing.T) {
	vacRepo := new(vacancyRepoMock)

	req := &create.Request{
		CompanyID:   uuid.New(),
		Title:       "Go Backend Dev",
		Description: "Go backend dev pretty much",
		WorkFormat:  vacancy.WorkFormatHybrid,
	}

	t.Run("global role wrong", func(t *testing.T) {
		memRepo := new(memRepoMock)
		uc := create.NewUsecase(vacRepo, memRepo)

		_, err := uc.Execute(
			context.Background(),
			req,
			&identity.Identity{UserID: uuid.New(), Role: identity.RoleCandidate},
		)

		require.ErrorIs(t, err, identity.ErrHrRoleRequired)
		memRepo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
		vacRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("hr is not member of company", func(t *testing.T) {
		memRepo := new(memRepoMock)
		ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

		memRepo.On("Get", mock.Anything, ident.UserID, req.CompanyID).
			Return(nil, member.ErrCompanyMemberNotFound).Once()

		uc := create.NewUsecase(vacRepo, memRepo)

		_, err := uc.Execute(context.Background(), req, ident)

		require.ErrorIs(t, err, member.ErrCompanyMemberRequired)
		memRepo.AssertExpectations(t)
		vacRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})
}

func TestUsecase_Execute_VacCreateFail(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	memRepo := new(memRepoMock)

	req := &create.Request{
		CompanyID:   uuid.New(),
		Title:       "Go Backend Dev",
		Description: "Go backend dev pretty much",
		WorkFormat:  vacancy.WorkFormatHybrid,
	}

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	memRepo.On("Get", mock.Anything, ident.UserID, req.CompanyID).
		Return(&member.CompanyMember{Role: member.CompanyRoleRecruiter}, nil).Once()

	vacRepo.On("Create", mock.Anything, mock.Anything).
		Return(errors.New("some err")).Once()

	uc := create.NewUsecase(vacRepo, memRepo)

	_, err := uc.Execute(context.Background(), req, ident)

	require.EqualError(t, err, "some err")
	memRepo.AssertExpectations(t)
	vacRepo.AssertExpectations(t)
}

func TestUsecase_Execute_ValidateErr(t *testing.T) {
	vacRepo := new(vacancyRepoMock)
	memRepo := new(memRepoMock)

	invalidReq := &create.Request{
		CompanyID: uuid.New(),
	}

	uc := create.NewUsecase(vacRepo, memRepo)

	ident := &identity.Identity{UserID: uuid.New(), Role: identity.RoleHR}

	_, err := uc.Execute(context.Background(), invalidReq, ident)

	require.Error(t, err)
	memRepo.AssertNotCalled(t, "Get", mock.Anything, mock.Anything, mock.Anything)
	vacRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}
