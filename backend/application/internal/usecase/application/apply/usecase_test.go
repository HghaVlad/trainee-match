package apply_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/apply"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/apply/mocks"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/identity"
)

type fakeTxManager struct{}

func (f *fakeTxManager) Do(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	return fn(ctx)
}

type fakeSnapIDGetter struct{}

func (f fakeSnapIDGetter) GetDeterministicAppSnapshotID(
	_ projection.ResumeData,
	_ projection.Candidate,
) (uuid.UUID, error) {
	return uuid.New(), nil
}

type tDeps struct {
	candRepo    *mocks.MockcandidateProjRepo
	resumeRepo  *mocks.MockresumeProjRepo
	vacRepo     *mocks.MockvacProjRepo
	appRepo     *mocks.MockappRepo
	appSnapRepo *mocks.MockappSnapshotRepo
	historyRepo *mocks.MockappHistoryRepo
}

func setup(t *testing.T) *tDeps {
	ctrl := gomock.NewController(t)

	return &tDeps{
		candRepo:    mocks.NewMockcandidateProjRepo(ctrl),
		resumeRepo:  mocks.NewMockresumeProjRepo(ctrl),
		vacRepo:     mocks.NewMockvacProjRepo(ctrl),
		appRepo:     mocks.NewMockappRepo(ctrl),
		appSnapRepo: mocks.NewMockappSnapshotRepo(ctrl),
		historyRepo: mocks.NewMockappHistoryRepo(ctrl),
	}
}

func newUc(deps *tDeps) *apply.Usecase {
	return apply.NewUsecase(
		deps.appRepo,
		deps.resumeRepo,
		deps.candRepo,
		deps.vacRepo,
		deps.appSnapRepo,
		deps.historyRepo,
		new(fakeSnapIDGetter),
		new(fakeTxManager),
	)
}

func candidateIdentity() identity.Identity {
	return identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleCandidate,
	}
}

func TestUsecase_Execute_Success(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUc(deps)

	ident := candidateIdentity()

	candID := uuid.New()
	resumeID := uuid.New()
	vacID := uuid.New()
	compID := uuid.New()
	coverLetter := "hello"

	req := apply.Request{
		VacancyID:   vacID,
		ResumeID:    resumeID,
		CoverLetter: &coverLetter,
	}

	candidate := &projection.Candidate{
		ID:       candID,
		FullName: "Arsen",
		Email:    "arsen@test.com",
	}

	resume := &projection.Resume{
		ID:          resumeID,
		CandidateID: candID,
		Name:        "Go Resume",
		Status:      projection.ResumeStatusPublished,
		Data: projection.ResumeData{
			Email: "resume@test.com",
		},
	}

	vacancy := &projection.Vacancy{
		ID:            vacID,
		CompanyID:     compID,
		CompanyName:   "Acme",
		Title:         "Go Backend",
		Status:        projection.VacancyStatusPublished,
		ModStatus:     projection.ModerationStatusOK,
		CompModStatus: projection.ModerationStatusOK,
	}

	deps.candRepo.EXPECT().
		GetByUserID(gomock.Any(), ident.UserID).
		Return(candidate, nil)

	deps.resumeRepo.EXPECT().
		GetByID(gomock.Any(), resumeID).
		Return(resume, nil)

	deps.vacRepo.EXPECT().
		GetByID(gomock.Any(), vacID).
		Return(vacancy, nil)

	deps.appSnapRepo.EXPECT().
		CreateIdempotent(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, snap application.Snapshot) error {
			require.Equal(t, resumeID, snap.ResumeID)
			require.Equal(t, candID, snap.CandidateID)
			require.Equal(t, "Arsen", snap.FullName)

			return nil
		})

	deps.appRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, app application.Application) error {
			require.Equal(t, candID, app.CandidateID)
			require.Equal(t, vacID, app.VacancyID)
			require.Equal(t, compID, app.CompanyID)
			require.Equal(t, application.StatusSubmitted, app.Status)

			return nil
		})

	deps.historyRepo.EXPECT().
		Add(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, status application.StatusChange) error {
			require.Equal(t, application.StatusSubmitted, status.Status)
			require.Equal(t, application.ActorCandidate, status.ChangedByRole)

			return nil
		})

	resp, err := uc.Execute(context.Background(), req, ident)

	require.NoError(t, err)
	require.NotNil(t, resp)

	require.Equal(t, vacID, resp.VacancyID)
	require.Equal(t, compID, resp.CompanyID)
	require.Equal(t, "Go Backend", resp.VacancyTitle)
	require.Equal(t, application.StatusSubmitted, resp.Status)
	require.Len(t, resp.StatusHistory, 1)
	require.Len(t, resp.AllowedActions, 1)
}

func TestUsecase_Execute_NotCandidate(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUc(deps)

	req := apply.Request{}

	ident := identity.Identity{
		UserID: uuid.New(),
		Role:   identity.RoleHR,
	}

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, identity.ErrCandidateRoleRequired)
	require.Nil(t, resp)
}

func TestUsecase_Execute_CandidateNotFound(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUc(deps)

	ident := candidateIdentity()

	req := apply.Request{}

	deps.candRepo.EXPECT().
		GetByUserID(gomock.Any(), ident.UserID).
		Return(nil, projection.ErrCandidateNotFound)

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, projection.ErrCandidateNotFound)
	require.Nil(t, resp)
}

func TestUsecase_Execute_ResumeNotPublished(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUc(deps)

	ident := candidateIdentity()

	candID := uuid.New()
	resumeID := uuid.New()

	req := apply.Request{
		ResumeID: resumeID,
	}

	deps.candRepo.EXPECT().
		GetByUserID(gomock.Any(), ident.UserID).
		Return(&projection.Candidate{
			ID: candID,
		}, nil)

	deps.resumeRepo.EXPECT().
		GetByID(gomock.Any(), resumeID).
		Return(&projection.Resume{
			ID:          resumeID,
			CandidateID: candID,
			Status:      projection.ResumeStatusDraft,
		}, nil)

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, application.ErrResumeNotPublished)
	require.Nil(t, resp)
}

func TestUsecase_Execute_ResumeAccessDenied(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUc(deps)

	ident := candidateIdentity()

	req := apply.Request{
		ResumeID: uuid.New(),
	}

	deps.candRepo.EXPECT().
		GetByUserID(gomock.Any(), ident.UserID).
		Return(&projection.Candidate{
			ID: uuid.New(),
		}, nil)

	deps.resumeRepo.EXPECT().
		GetByID(gomock.Any(), req.ResumeID).
		Return(&projection.Resume{
			ID:          req.ResumeID,
			CandidateID: uuid.New(),
			Status:      projection.ResumeStatusPublished,
		}, nil)

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, application.ErrResumeAccessDenied)
	require.Nil(t, resp)
}

func TestUsecase_Execute_VacancyNotApplyable(t *testing.T) {
	t.Parallel()

	ident := candidateIdentity()

	candID := uuid.New()

	req := apply.Request{
		ResumeID:  uuid.New(),
		VacancyID: uuid.New(),
	}

	tts := []struct {
		vac *projection.Vacancy
		expectedErr error
	}{
		{
			vac: &projection.Vacancy{
				ID:            req.VacancyID,
				Status:        projection.VacancyStatusArchived,
				ModStatus:     projection.ModerationStatusOK,
				CompModStatus: projection.ModerationStatusOK,
			},
			expectedErr: projection.ErrVacancyNotPublished,
		},
		{
			vac: &projection.Vacancy{
				ID:            req.VacancyID,
				Status:        projection.VacancyStatusPublished,
				ModStatus:     projection.ModerationStatusHidden,
				CompModStatus: projection.ModerationStatusOK,
			},
			expectedErr: projection.ErrVacancyBadModStatus,
		},
		{
			vac: &projection.Vacancy{
				ID:            req.VacancyID,
				Status:        projection.VacancyStatusPublished,
				ModStatus:     projection.ModerationStatusOK,
				CompModStatus: projection.ModerationStatusHidden,
			},
			expectedErr: projection.ErrVacancyBadModStatus,
		},
	}

	for _, tt := range tts {
		deps := setup(t)
		uc := newUc(deps)

		deps.candRepo.EXPECT().
			GetByUserID(gomock.Any(), ident.UserID).
			Return(&projection.Candidate{
				ID: candID,
			}, nil)

		deps.resumeRepo.EXPECT().
			GetByID(gomock.Any(), req.ResumeID).
			Return(&projection.Resume{
				ID:          req.ResumeID,
				CandidateID: candID,
				Status:      projection.ResumeStatusPublished,
			}, nil)

		deps.vacRepo.EXPECT().
			GetByID(gomock.Any(), req.VacancyID).
			Return(tt.vac, nil)

		resp, err := uc.Execute(context.Background(), req, ident)

		require.ErrorIs(t, tt.expectedErr, err)
		require.Nil(t, resp)
	}
}

func TestUsecase_Execute_CreateApplicationError(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUc(deps)

	ident := candidateIdentity()

	candID := uuid.New()
	resumeID := uuid.New()
	vacID := uuid.New()

	req := apply.Request{
		ResumeID:  resumeID,
		VacancyID: vacID,
	}

	deps.candRepo.EXPECT().
		GetByUserID(gomock.Any(), ident.UserID).
		Return(&projection.Candidate{
			ID:       candID,
			FullName: "Arsen",
			Email:    "test@test.com",
		}, nil)

	deps.resumeRepo.EXPECT().
		GetByID(gomock.Any(), resumeID).
		Return(&projection.Resume{
			ID:          resumeID,
			CandidateID: candID,
			Status:      projection.ResumeStatusPublished,
		}, nil)

	deps.vacRepo.EXPECT().
		GetByID(gomock.Any(), vacID).
		Return(&projection.Vacancy{
			ID:            vacID,
			CompanyID:     uuid.New(),
			Status:        projection.VacancyStatusPublished,
			ModStatus:     projection.ModerationStatusOK,
			CompModStatus: projection.ModerationStatusOK,
		}, nil)

	deps.appSnapRepo.EXPECT().
		CreateIdempotent(gomock.Any(), gomock.Any()).
		Return(nil)

	expectedErr := application.ErrActiveAlreadyExists

	deps.appRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(expectedErr)

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, resp)
}

func TestUsecase_Execute_AddHistoryError(t *testing.T) {
	t.Parallel()

	deps := setup(t)
	uc := newUc(deps)

	ident := candidateIdentity()

	candID := uuid.New()
	resumeID := uuid.New()
	vacID := uuid.New()

	req := apply.Request{
		ResumeID:  resumeID,
		VacancyID: vacID,
	}

	deps.candRepo.EXPECT().
		GetByUserID(gomock.Any(), ident.UserID).
		Return(&projection.Candidate{
			ID:       candID,
			FullName: "Arsen",
			Email:    "test@test.com",
		}, nil)

	deps.resumeRepo.EXPECT().
		GetByID(gomock.Any(), resumeID).
		Return(&projection.Resume{
			ID:          resumeID,
			CandidateID: candID,
			Status:      projection.ResumeStatusPublished,
		}, nil)

	deps.vacRepo.EXPECT().
		GetByID(gomock.Any(), vacID).
		Return(&projection.Vacancy{
			ID:            vacID,
			CompanyID:     uuid.New(),
			Status:        projection.VacancyStatusPublished,
			ModStatus:     projection.ModerationStatusOK,
			CompModStatus: projection.ModerationStatusOK,
		}, nil)

	deps.appSnapRepo.EXPECT().
		CreateIdempotent(gomock.Any(), gomock.Any()).
		Return(nil)

	deps.appRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)

	expectedErr := errors.New("history error")

	deps.historyRepo.EXPECT().
		Add(gomock.Any(), gomock.Any()).
		Return(expectedErr)

	resp, err := uc.Execute(context.Background(), req, ident)

	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, resp)
}
