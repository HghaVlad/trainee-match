package create_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/create"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/create/mocks"
)

type fakeTxManager struct{}

func (f *fakeTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestUsecase_Execute_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	compRepo := mocks.NewMockCompanyRepo(ctrl)
	memberRepo := mocks.NewMockCompanyMemberRepo(ctrl)
	outboxWriter := mocks.NewMockoutboxWriter(ctrl)

	txManager := &fakeTxManager{}

	uc := create.NewUsecase(
		compRepo,
		memberRepo,
		outboxWriter,
		txManager,
	)

	req := &create.Request{
		Name:        "Acme",
		Description: ptr("Best company"),
		Website:     ptr("https://acme.com"),
	}

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleHR,
	}

	compRepo.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, comp *company.Company) error {
			require.Equal(t, req.Name, comp.Name)
			require.Equal(t, req.Description, comp.Description)
			require.Equal(t, req.Website, comp.Website)
			require.NotEqual(t, uuid.Nil, comp.ID)

			return nil
		})

	memberRepo.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, memb *member.CompanyMember) error {
			require.Equal(t, ident.UserID, memb.UserID)
			require.Equal(t, member.CompanyRoleAdmin, memb.Role)
			require.NotEqual(t, uuid.Nil, memb.CompanyID)

			return nil
		})

	outboxWriter.
		EXPECT().
		WriteCompanyMemberAdded(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, ev member.AddedEvent) error {
			require.Equal(t, ident.UserID, ev.UserID)
			require.Equal(t, member.CompanyRoleAdmin, ev.Role)
			require.NotEqual(t, uuid.Nil, ev.CompanyID)
			require.NotEqual(t, uuid.Nil, ev.EventID)

			return nil
		})

	resp, err := uc.Execute(context.Background(), req, ident)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEqual(t, uuid.Nil, resp.ID)
}

func TestUsecase_Execute_HRRoleRequired(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	uc := create.NewUsecase(
		mocks.NewMockCompanyRepo(ctrl),
		mocks.NewMockCompanyMemberRepo(ctrl),
		mocks.NewMockoutboxWriter(ctrl),
		&fakeTxManager{},
	)

	req := &create.Request{
		Name:        "Acme",
		Description: ptr("Best company"),
		Website:     ptr("https://acme.com"),
	}

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleCandidate,
	}

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, identity.ErrHrRoleRequired)
	require.Nil(t, resp)
}

func TestUsecase_Execute_ValidationError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	uc := create.NewUsecase(
		mocks.NewMockCompanyRepo(ctrl),
		mocks.NewMockCompanyMemberRepo(ctrl),
		mocks.NewMockoutboxWriter(ctrl),
		&fakeTxManager{},
	)

	req := &create.Request{
		Name: "",
	}

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleHR,
	}

	resp, err := uc.Execute(context.Background(), req, ident)

	require.Error(t, err)
	require.Nil(t, resp)
}

func TestUsecase_Execute_AlreadyExists(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	compRepo := mocks.NewMockCompanyRepo(ctrl)

	uc := create.NewUsecase(
		compRepo,
		mocks.NewMockCompanyMemberRepo(ctrl),
		mocks.NewMockoutboxWriter(ctrl),
		&fakeTxManager{},
	)

	req := &create.Request{
		Name:        "Acme",
		Description: ptr("Best company"),
		Website:     ptr("https://acme.com"),
	}

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleHR,
	}

	compRepo.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(company.ErrCompanyAlreadyExists)

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, company.ErrCompanyAlreadyExists)
	require.Nil(t, resp)
}

func TestUsecase_Execute_MemberCreateError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	compRepo := mocks.NewMockCompanyRepo(ctrl)
	memberRepo := mocks.NewMockCompanyMemberRepo(ctrl)

	uc := create.NewUsecase(
		compRepo,
		memberRepo,
		mocks.NewMockoutboxWriter(ctrl),
		&fakeTxManager{},
	)

	req := &create.Request{
		Name:        "Acme",
		Description: ptr("Best company"),
		Website:     ptr("https://acme.com"),
	}

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleHR,
	}

	expectedErr := errors.New("create member failed")

	compRepo.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	memberRepo.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(expectedErr)

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, resp)
}

func TestUsecase_Execute_OutboxError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	compRepo := mocks.NewMockCompanyRepo(ctrl)
	memberRepo := mocks.NewMockCompanyMemberRepo(ctrl)
	outboxWriter := mocks.NewMockoutboxWriter(ctrl)

	uc := create.NewUsecase(
		compRepo,
		memberRepo,
		outboxWriter,
		&fakeTxManager{},
	)

	req := &create.Request{
		Name:        "Acme",
		Description: ptr("Best company"),
		Website:     ptr("https://acme.com"),
	}

	ident := &identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleHR,
	}

	expectedErr := errors.New("outbox failed")

	compRepo.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	memberRepo.
		EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	outboxWriter.
		EXPECT().
		WriteCompanyMemberAdded(gomock.Any(), gomock.Any()).
		Return(expectedErr)

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, resp)
}

func ptr[T any](val T) *T {
	return &val
}
