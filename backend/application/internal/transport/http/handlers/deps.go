package handlers

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/analytics/dynamics"

	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/analytics/summary"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/apply"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/hrupdstatus"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/listcandidatesummary"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/listhrsummary"
	appviews "github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/withdraw"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/identity"
)

type Deps struct {
	Apply               applyUC
	ListCandidateApps   listCandidateApps
	GetCandidateViewUC  getCandidateDetailedView
	GetCandiStatHistory getCandidateStatusHistory
	Withdraw            withdrawUC
	ListHrApps          listHrApps
	GetHrDetailedView   getHrDetailedView
	HrUpdateStatus      hrUpdateStatus
	GetHistoryHrView    getHistoryHrView
	AnalyticsSummary    analyticsSummary
	DynamicsDashboard   dynamicsDashboard
	Logger              *slog.Logger
}

type applyUC interface {
	Execute(context.Context, apply.Request, identity.Identity) (*appviews.CandidateDetailedView, error)
}

type getCandidateDetailedView interface {
	Execute(ctx context.Context, appID uuid.UUID, ident identity.Identity) (*appviews.CandidateDetailedView, error)
}

type getCandidateStatusHistory interface {
	Execute(
		ctx context.Context,
		appID uuid.UUID,
		ident identity.Identity,
	) ([]appviews.StatusChangeCandidateFullView, error)
}

type listCandidateApps interface {
	Execute(
		ctx context.Context,
		req listcandidatesummary.Request,
		ident identity.Identity,
	) (*listcandidatesummary.Response, error)
}

type withdrawUC interface {
	Execute(
		ctx context.Context,
		req withdraw.Request,
		ident identity.Identity,
	) (*appviews.CandidateDetailedView, error)
}

type listHrApps interface {
	Execute(
		ctx context.Context,
		req listhrsummary.Request,
		ident identity.Identity,
	) (*listhrsummary.Response, error)
}

type getHrDetailedView interface {
	Execute(
		ctx context.Context,
		appID uuid.UUID,
		ident identity.Identity,
	) (*appviews.HrDetailedView, error)
}

type hrUpdateStatus interface {
	Execute(
		ctx context.Context,
		req hrupdstatus.Request,
		ident identity.Identity,
	) (*appviews.HrDetailedView, error)
}

type getHistoryHrView interface {
	Execute(
		ctx context.Context,
		appID uuid.UUID,
		ident identity.Identity,
	) ([]appviews.StatusChangeHrFullView, error)
}

type analyticsSummary interface {
	GetByCompany(
		ctx context.Context,
		compID uuid.UUID,
		ident identity.Identity,
	) (*summary.Summary, error)

	GetByVacancy(
		ctx context.Context,
		compID uuid.UUID,
		ident identity.Identity,
	) (*summary.Summary, error)
}

type dynamicsDashboard interface {
	GetDashboardByCompany(
		ctx context.Context,
		compID uuid.UUID,
		period dynamics.Period,
		iden identity.Identity,
	) ([]dynamics.Bucket, error)

	GetDashboardByVacancy(
		ctx context.Context,
		vacID uuid.UUID,
		period dynamics.Period,
		iden identity.Identity,
	) ([]dynamics.Bucket, error)
}
