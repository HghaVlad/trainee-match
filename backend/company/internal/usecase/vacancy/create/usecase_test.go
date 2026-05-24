package create_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/create"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/create/mocks"
)

type fakeTxManager struct {
	called bool
}

func (f *fakeTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	f.called = true
	return fn(ctx)
}

type testDeps struct {
	vacRepo      *mocks.MockVacancyRepo
	compRepo     *mocks.MockCompanyRepo
	outboxWriter *mocks.MockoutboxWriter
}

func setup(t *testing.T) *testDeps {
	ctrl := gomock.NewController(t)

	return &testDeps{
		vacRepo:      mocks.NewMockVacancyRepo(ctrl),
		compRepo:     mocks.NewMockCompanyRepo(ctrl),
		outboxWriter: mocks.NewMockoutboxWriter(ctrl),
	}
}

func newUC(deps *testDeps) *create.Usecase {
	return create.NewUsecase(
		deps.vacRepo,
		deps.compRepo,
		deps.outboxWriter,
		new(fakeTxManager),
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

	deps.outboxWriter.EXPECT().
		WriteVacancyDraftCreated(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, ev vacancy.DraftCreatedEvent) error {
			require.Equal(t, companyID, ev.CompanyID)
			require.Equal(t, req.Title, ev.Title)
			require.Equal(t, "Acme", ev.CompanyName)
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
