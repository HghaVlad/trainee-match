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

func (u *Usecase) Execute(
	ctx context.Context,
	req Request,
	ident identity.Identity,
) (*views.CandidateDetailedView, error) {
	if ident.Role != identity.RoleCandidate {
		return nil, identity.ErrCandidateRoleRequired
	}

	if err := req.validate(); err != nil {
		return nil, err
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

	vacProj, err := u.getApplyableVacancy(ctx, req.VacancyID)
	if err != nil {
		return nil, err
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

	return buildFullCandidateView(app, appSnap, vacProj), nil
}

func (u *Usecase) getResumeProjAndCheck(ctx context.Context, resID, candID uuid.UUID) (*projection.Resume, error) {
	resumeProj, err := u.resumeProjRepo.GetByID(ctx, resID)
	if err != nil {
		return nil, err
	}

	if resumeProj.Status != projection.ResumeStatusPublished {
		return nil, application.ErrResumeNotPublished
	}
	if resumeProj.ModerationStatus != projection.ModerationStatusOK {
		return nil, projection.ErrResumeBadModStatus
	}

	if resumeProj.CandidateID != candID {
		return nil, application.ErrResumeAccessDenied
	}

	return resumeProj, nil
}

func (u *Usecase) getApplyableVacancy(ctx context.Context, vacID uuid.UUID) (*projection.Vacancy, error) {
	vacProj, err := u.vacProjRepo.GetByID(ctx, vacID)
	if err != nil {
		return nil, err
	}

	if err := vacProj.IsApplyable(); err != nil {
		return nil, err
	}

	return vacProj, nil
}

func buildFullCandidateView(
	app *application.Application,
	appSnap *application.Snapshot,
	vacProj *projection.Vacancy,
) *views.CandidateDetailedView {
	return &views.CandidateDetailedView{
		AppID:        app.ID,
		VacancyID:    vacProj.ID,
		CompanyID:    vacProj.CompanyID,
		VacancyTitle: vacProj.Title,
		CompanyName:  vacProj.CompanyName,
		Status:       app.Status,
		CoverLetter:  app.CoverLetter,
		Snapshot: views.ApplicationSnapshot{
			ResumeData: appSnap.ResumeData,
			Email:      appSnap.Email,
			FullName:   appSnap.FullName,
			Telegram:   appSnap.Telegram,
			CreatedAt:  appSnap.CreatedAt,
		},
		CreatedAt: app.CreatedAt,
		UpdatedAt: app.UpdatedAt,
		StatusHistory: []views.StatusChangeCandidateView{
			{
				Status:        app.Status,
				ChangedByRole: application.ActorCandidate,
				CreatedAt:     app.CreatedAt},
		},
		AllowedActions: []views.AllowedAction{views.AllowedActionWithdraw},
	}
}
