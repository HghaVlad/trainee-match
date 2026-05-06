package apply

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/identity"
)

type Usecase struct {
	appRepo        appRepo
	resumeProjRepo resumeProjRepo
	candProjRepo   candidateProjRepo
	vacProjRepo    vacProjRepo
	appSnapRepo    appSnapshotRepo
	appHistoryRepo appHistoryRepo
	snapIDGetter   snapIDGetter
	txManager      common.TxManager
}

func NewUsecase(
	appRepo appRepo,
	resProjRepo resumeProjRepo,
	candProjRepo candidateProjRepo,
	vacProjRepo vacProjRepo,
	appSnapRepo appSnapshotRepo,
	appHistoryRepo appHistoryRepo,
	snapIDGetter snapIDGetter,
	txManager common.TxManager,
) *Usecase {
	return &Usecase{
		appRepo:        appRepo,
		resumeProjRepo: resProjRepo,
		candProjRepo:   candProjRepo,
		vacProjRepo:    vacProjRepo,
		appSnapRepo:    appSnapRepo,
		appHistoryRepo: appHistoryRepo,
		snapIDGetter:   snapIDGetter,
		txManager:      txManager,
	}
}

func (u *Usecase) Execute(ctx context.Context, req Request, ident identity.Identity) (*views.Details, error) {
	if ident.Role != identity.RoleCandidate {
		return nil, identity.ErrCandidateRoleRequired
	}

	now := time.Now().UTC()

	candProj, err := u.candProjRepo.GetByUserID(ctx, ident.UserID)
	if err != nil {
		return nil, err
	}

	resumeProj, err := u.getResumeProjAndCheck(ctx, req.ResumeID, candProj.ID)
	if err != nil {
		return nil, err
	}

	vacProj, err := u.vacProjRepo.GetByID(ctx, req.VacancyID)
	if err != nil {
		return nil, err
	}

	if vacProj.Status != projection.VacancyStatusPublished {
		return nil, application.ErrVacancyNotPublished
	}

	// for deduplication of snapshots
	appSnapID, err := u.snapIDGetter.GetDeterministicAppSnapshotID(resumeProj.Data, *candProj)
	if err != nil {
		return nil, err
	}

	appSnap, err := application.NewApplicationSnapshot(appSnapID, *resumeProj, *candProj, now)
	if err != nil {
		return nil, err
	}

	app, err := application.NewSubmitted(
		resumeProj.ID, candProj.ID,
		vacProj.ID, vacProj.CompanyID,
		appSnap.ID, req.CoverLetter,
		now,
	)
	if err != nil {
		return nil, err
	}

	statusChange, err := application.NewStatusChange(
		app.ID, application.StatusSubmitted,
		&ident.UserID, application.ActorCandidate,
		nil, now,
	)
	if err != nil {
		return nil, err
	}

	err = u.txManager.Do(ctx, func(ctx context.Context) error {
		err = u.appSnapRepo.CreateIdempotent(ctx, *appSnap)
		if err != nil {
			return err
		}

		err = u.appRepo.Create(ctx, *app)
		if err != nil {
			return err
		}

		return u.appHistoryRepo.Add(ctx, *statusChange)
	})
	if err != nil {
		return nil, err
	}

	return &views.Details{
		Application:    app,
		VacProj:        vacProj,
		Snapshot:       appSnap,
		AllowedActions: []views.AllowedAction{views.AllowedActionWithdraw},
	}, nil
}

func (u *Usecase) getResumeProjAndCheck(ctx context.Context, resID, candID uuid.UUID) (*projection.Resume, error) {
	resumeProj, err := u.resumeProjRepo.GetByID(ctx, resID)
	if err != nil {
		return nil, err
	}

	if resumeProj.Status != projection.ResumeStatusPublished {
		return nil, application.ErrResumeNotPublished
	}

	if resumeProj.CandidateID != candID {
		return nil, application.ErrResumeAccessDenied
	}

	return resumeProj, nil
}
