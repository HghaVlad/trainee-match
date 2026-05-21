package create_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/create"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/create/mocks"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/views"
)

type testDeps struct {
	vacRepo    *mocks.MockVacancyRepo
	memberRepo *mocks.MockCompMemberRepo
	compRepo   *mocks.MockCompanyRepo
	searchRepo *mocks.MockSearchRepo
}

func setup(t *testing.T) *testDeps {
	ctrl := gomock.NewController(t)

	return &testDeps{
		vacRepo:    mocks.NewMockVacancyRepo(ctrl),
		memberRepo: mocks.NewMockCompMemberRepo(ctrl),
		compRepo:   mocks.NewMockCompanyRepo(ctrl),
		searchRepo: mocks.NewMockSearchRepo(ctrl),
	}
}

func newUC(deps *testDeps) *create.Usecase {
	return create.NewUsecase(
		deps.vacRepo,
		deps.compRepo,
		deps.searchRepo,
	)
}

func TestUsecase_Execute_Success(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUC(deps)

	companyID := uuid.New()
	userID := uuid.New()

	req := &create.Request{
		CompanyID:   companyID,
		Title:       "Backend Go Developer",
		Description: "Develop scalable backend services",

		WorkFormat: vacancy.WorkFormatRemote,

		IsPaid: true,
	}

	ident := &identity.Identity{
		UserID: userID,
	}

	comp := &company.Company{
		ID:   companyID,
		Name: "Acme",
	}

	deps.compRepo.EXPECT().
		GetByMember(gomock.Any(), companyID, userID).
		Return(comp, nil)

	deps.vacRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, vac *vacancy.Vacancy) error {
			require.Equal(t, companyID, vac.CompanyID)
			require.Equal(t, userID, vac.CreatedBy)
			require.Equal(t, vacancy.StatusDraft, vac.Status)
			require.Equal(t, vacancy.ModerationStatusOK, vac.ModerationStatus)

			require.Equal(t, req.Title, vac.Title)
			require.Equal(t, req.Description, vac.Description)

			require.NotEqual(t, uuid.Nil, vac.ID)

			return nil
		})

	deps.searchRepo.EXPECT().
		Index(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, search views.VacancySearch) error {
			require.Equal(t, companyID, search.CompanyID)
			require.Equal(t, "Acme", search.CompanyName)

			require.Equal(t, req.Title, search.Title)
			require.Equal(t, req.Description, search.Description)

			require.Equal(t, vacancy.StatusDraft, search.Status)

			return nil
		})

	resp, err := uc.Execute(context.Background(), req, ident)
	require.NoError(t, err)

	require.NotNil(t, resp)
	require.NotEqual(t, uuid.Nil, resp.ID)
}

func TestUsecase_Execute_ValidateError(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUC(deps)

	req := &create.Request{}

	ident := &identity.Identity{
		UserID: uuid.New(),
	}

	resp, err := uc.Execute(context.Background(), req, ident)

	require.Error(t, err)
	require.Nil(t, resp)
}

func TestUsecase_Execute_GetCompanyByMemberError(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUC(deps)

	companyID := uuid.New()
	userID := uuid.New()

	req := &create.Request{
		CompanyID:   companyID,
		Title:       "Backend Developer",
		Description: "Description",
		WorkFormat:  vacancy.WorkFormatRemote,
	}

	ident := &identity.Identity{
		UserID: userID,
	}

	deps.compRepo.EXPECT().
		GetByMember(gomock.Any(), companyID, userID).
		Return(nil, company.ErrCompanyNotFound)

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, company.ErrCompanyNotFound)
	require.Nil(t, resp)
}

func TestUsecase_Execute_CreateVacancyError(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUC(deps)

	companyID := uuid.New()
	userID := uuid.New()

	req := &create.Request{
		CompanyID:   companyID,
		Title:       "Backend Developer",
		Description: "Description",
		WorkFormat:  vacancy.WorkFormatRemote,
	}

	ident := &identity.Identity{
		UserID: userID,
	}

	comp := &company.Company{
		ID:   companyID,
		Name: "Acme",
	}

	expectedErr := errors.New("create vacancy error")

	deps.compRepo.EXPECT().
		GetByMember(gomock.Any(), companyID, userID).
		Return(comp, nil)

	deps.vacRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(expectedErr)

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, resp)
}

func TestUsecase_Execute_IndexSearchError(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUC(deps)

	companyID := uuid.New()
	userID := uuid.New()

	req := &create.Request{
		CompanyID:   companyID,
		Title:       "Backend Developer",
		Description: "Description",
		WorkFormat:  vacancy.WorkFormatRemote,
	}

	ident := &identity.Identity{
		UserID: userID,
	}

	comp := &company.Company{
		ID:   companyID,
		Name: "Acme",
	}

	expectedErr := errors.New("index error")

	deps.compRepo.EXPECT().
		GetByMember(gomock.Any(), companyID, userID).
		Return(comp, nil)

	deps.vacRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	deps.searchRepo.EXPECT().
		Index(gomock.Any(), gomock.Any()).
		Return(expectedErr)

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, resp)
}
