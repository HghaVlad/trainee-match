package handlers

import (
	"context"
	"errors"
	"log/slog"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/projection"
	"github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/mappers"
	"github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/middleware"
	"github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/oapi"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/analytics/dynamics"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/hrupdstatus"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/listhrsummary"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/withdraw"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/cursors"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/identity"
)

type Handler struct {
	apply               applyUC
	listCandidateApps   listCandidateApps
	getCandidateView    getCandidateDetailedView
	getHistoryCandiView getCandidateStatusHistory
	withdraw            withdrawUC
	listHrApps          listHrApps
	getHrDetailedView   getHrDetailedView
	hrUpdateStatus      hrUpdateStatus
	getHistoryHrView    getHistoryHrView
	analyticsSummary    analyticsSummary
	dynamicsDashboard   dynamicsDashboard
	logger              *slog.Logger
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		apply:               deps.Apply,
		listCandidateApps:   deps.ListCandidateApps,
		getCandidateView:    deps.GetCandidateViewUC,
		getHistoryCandiView: deps.GetCandiStatHistory,
		listHrApps:          deps.ListHrApps,
		withdraw:            deps.Withdraw,
		getHrDetailedView:   deps.GetHrDetailedView,
		hrUpdateStatus:      deps.HrUpdateStatus,
		getHistoryHrView:    deps.GetHistoryHrView,
		analyticsSummary:    deps.AnalyticsSummary,
		dynamicsDashboard:   deps.DynamicsDashboard,
		logger:              deps.Logger,
	}
}

func (h *Handler) CreateApplication(
	ctx context.Context,
	request oapi.CreateApplicationRequestObject,
) (oapi.CreateApplicationResponseObject, error) {
	ident := middleware.IdentityFromContext(ctx)
	req := mappers.ApplyReqToUC(request)

	resp, err := h.apply.Execute(ctx, req, *ident)

	if err != nil {
		return applyErrToResponse(err)
	}

	return oapi.CreateApplication201JSONResponse{
		Data: mappers.CandidateDetailedViewToHTTP(resp),
	}, nil
}

func (h *Handler) GetMyApplication(
	ctx context.Context,
	request oapi.GetMyApplicationRequestObject,
) (oapi.GetMyApplicationResponseObject, error) {
	ident := middleware.IdentityFromContext(ctx)
	appID := request.ApplicationId

	resp, err := h.getCandidateView.Execute(ctx, appID, *ident)
	if err != nil {
		return handleCandViewErr(err)
	}

	return oapi.GetMyApplication200JSONResponse{
		Data: mappers.CandidateDetailedViewToHTTP(resp),
	}, nil
}

func (h *Handler) GetMyApplicationHistory(
	ctx context.Context,
	request oapi.GetMyApplicationHistoryRequestObject,
) (oapi.GetMyApplicationHistoryResponseObject, error) {
	ident := middleware.IdentityFromContext(ctx)
	appID := request.ApplicationId

	history, err := h.getHistoryCandiView.Execute(ctx, appID, *ident)
	if err != nil {
		return handleCandiHistoryErr(err)
	}

	return mappers.CandidateHistoryFullToHTTP(history), nil
}

func (h *Handler) ListMyApplications(
	ctx context.Context,
	request oapi.ListMyApplicationsRequestObject,
) (oapi.ListMyApplicationsResponseObject, error) {
	ident := middleware.IdentityFromContext(ctx)
	req := mappers.ListMyApplicationsReqToUC(request)

	resp, err := h.listCandidateApps.Execute(ctx, req, *ident)
	if err != nil {
		return handleListMyAppErr(err)
	}

	return mappers.CandidateListResponseToHTTP(resp), nil
}

func (h *Handler) WithdrawApplication(
	ctx context.Context,
	request oapi.WithdrawApplicationRequestObject,
) (oapi.WithdrawApplicationResponseObject, error) {
	ident := middleware.IdentityFromContext(ctx)
	req := withdraw.Request{
		AppID:   request.ApplicationId,
		Comment: request.Body.Comment,
	}

	view, err := h.withdraw.Execute(ctx, req, *ident)
	if err != nil {
		return handleWithdrawErr(err)
	}

	return oapi.WithdrawApplication200JSONResponse{
		Data: mappers.CandidateDetailedViewToHTTP(view),
	}, nil
}

func (h *Handler) ListCompanyApplications(
	ctx context.Context,
	request oapi.ListCompanyApplicationsRequestObject,
) (oapi.ListCompanyApplicationsResponseObject, error) {
	ident := middleware.IdentityFromContext(ctx)
	req := mappers.ListCompanyApplicationsReqToUC(request)

	resp, err := h.listHrApps.Execute(ctx, req, *ident)

	if err != nil {
		return handleHrCompListErr(err)
	}

	return oapi.ListCompanyApplications200JSONResponse(mappers.HrListResponseToHTTP(resp)), nil
}

func (h *Handler) ListVacancyApplications(
	ctx context.Context,
	request oapi.ListVacancyApplicationsRequestObject,
) (oapi.ListVacancyApplicationsResponseObject, error) {
	ident := middleware.IdentityFromContext(ctx)
	req := mappers.ListVacancyApplicationsReqToUC(request)

	resp, err := h.listHrApps.Execute(ctx, req, *ident)

	if err != nil {
		return handleHrVacListErr(err)
	}

	return oapi.ListVacancyApplications200JSONResponse(mappers.HrListResponseToHTTP(resp)), nil
}

func (h *Handler) GetHrApplication(
	ctx context.Context,
	request oapi.GetHrApplicationRequestObject,
) (oapi.GetHrApplicationResponseObject, error) {
	ident := middleware.IdentityFromContext(ctx)
	appID := request.ApplicationId

	view, err := h.getHrDetailedView.Execute(ctx, appID, *ident)
	if err != nil {
		return handleHrViewErr(err)
	}

	return oapi.GetHrApplication200JSONResponse{
		Data: mappers.HrDetailedViewToHTTP(view),
	}, nil
}

func (h *Handler) ChangeApplicationStatus(
	ctx context.Context,
	request oapi.ChangeApplicationStatusRequestObject,
) (oapi.ChangeApplicationStatusResponseObject, error) {
	ident := middleware.IdentityFromContext(ctx)
	req := hrupdstatus.Request{
		AppID:   request.ApplicationId,
		Status:  application.Status(request.Body.Status),
		Comment: request.Body.Comment,
	}

	view, err := h.hrUpdateStatus.Execute(ctx, req, *ident)

	if err != nil {
		return handleHrUpdStatus(err)
	}

	return oapi.ChangeApplicationStatus200JSONResponse{
		Data: mappers.HrDetailedViewToHTTP(view),
	}, nil
}

func (h *Handler) GetHrApplicationHistory(
	ctx context.Context,
	request oapi.GetHrApplicationHistoryRequestObject,
) (oapi.GetHrApplicationHistoryResponseObject, error) {
	ident := middleware.IdentityFromContext(ctx)
	appID := request.ApplicationId

	history, err := h.getHistoryHrView.Execute(ctx, appID, *ident)
	if err != nil {
		return handleGetHistoryHrViewErr(err)
	}

	return oapi.GetHrApplicationHistory200JSONResponse{
		Data: mappers.HrHistoryToHTTP(history),
	}, nil
}

func (h *Handler) GetCompanyAnalyticsSummary(
	ctx context.Context,
	request oapi.GetCompanyAnalyticsSummaryRequestObject,
) (oapi.GetCompanyAnalyticsSummaryResponseObject, error) {
	ident := middleware.IdentityFromContext(ctx)
	compID := request.CompanyId

	sum, err := h.analyticsSummary.GetByCompany(ctx, compID, *ident)
	if err != nil {
		return handleCompanyAnalyticsSummary(err)
	}

	return oapi.GetCompanyAnalyticsSummary200JSONResponse{
		Data: mappers.AnalyticsSummaryToHTTP(*sum),
	}, nil
}

func (h *Handler) GetVacancyAnalyticsSummary(
	ctx context.Context,
	request oapi.GetVacancyAnalyticsSummaryRequestObject,
) (oapi.GetVacancyAnalyticsSummaryResponseObject, error) {
	ident := middleware.IdentityFromContext(ctx)
	vacID := request.VacancyId

	sum, err := h.analyticsSummary.GetByVacancy(ctx, vacID, *ident)
	if err != nil {
		return handleVacancyAnalyticsSummaryErr(err)
	}

	return oapi.GetVacancyAnalyticsSummary200JSONResponse{
		Data: mappers.AnalyticsSummaryToHTTP(*sum),
	}, nil
}

func (h *Handler) GetCompanyDynamics(
	ctx context.Context,
	request oapi.GetCompanyDynamicsRequestObject,
) (oapi.GetCompanyDynamicsResponseObject, error) {
	ident := middleware.IdentityFromContext(ctx)
	compID := request.CompanyId
	period := mappers.DynamicsCompPeriodToUC(request)

	dashboard, err := h.dynamicsDashboard.GetDashboardByCompany(ctx, compID, period, *ident)
	if err != nil {
		return handleCompanyDynamicsErr(err)
	}

	return oapi.GetCompanyDynamics200JSONResponse{
		Data: mappers.DynamicsDashboardToHTTP(dashboard),
	}, nil
}

func (h *Handler) GetVacancyDynamics(
	ctx context.Context,
	request oapi.GetVacancyDynamicsRequestObject,
) (oapi.GetVacancyDynamicsResponseObject, error) {
	ident := middleware.IdentityFromContext(ctx)
	vacID := request.VacancyId
	period := mappers.DynamicsVacPeriodToUC(request)

	dashboard, err := h.dynamicsDashboard.GetDashboardByVacancy(ctx, vacID, period, *ident)
	if err != nil {
		return handleVacancyDynamicsErr(err)
	}

	return oapi.GetVacancyDynamics200JSONResponse{
		Data: mappers.DynamicsDashboardToHTTP(dashboard),
	}, nil
}

func (h *Handler) GetCompanyStatusFunnel(
	ctx context.Context,
	request oapi.GetCompanyStatusFunnelRequestObject,
) (oapi.GetCompanyStatusFunnelResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) GetVacancyStatusFunnel(
	ctx context.Context,
	request oapi.GetVacancyStatusFunnelRequestObject,
) (oapi.GetVacancyStatusFunnelResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func applyErrToResponse(err error) (oapi.CreateApplicationResponseObject, error) {
	switch {
	case errors.Is(err, application.ErrCoverLetterTooLong):
		return oapi.CreateApplication400JSONResponse{
			BadRequestJSONResponse: oapi.BadRequestJSONResponse{
				Error:   "bad_request",
				Message: err.Error(),
			},
		}, nil

	case errors.Is(err, application.ErrResumeAccessDenied),
		errors.Is(err, identity.ErrCandidateRoleRequired):
		return oapi.CreateApplication403JSONResponse{
			ForbiddenErrorJSONResponse: oapi.ForbiddenErrorJSONResponse{
				Error:   "forbidden",
				Message: err.Error(),
			},
		}, nil

	case errors.Is(err, projection.ErrVacancyNotFound),
		errors.Is(err, projection.ErrResumeNotFound),
		errors.Is(err, projection.ErrCandidateNotFound):
		return oapi.CreateApplication404JSONResponse{
			Error:   "not_found",
			Message: err.Error(),
		}, nil

	case errors.Is(err, application.ErrVacancyNotPublished),
		errors.Is(err, application.ErrResumeNotPublished):
		return oapi.CreateApplication400JSONResponse{
			BadRequestJSONResponse: oapi.BadRequestJSONResponse{
				Error:   "bad_request",
				Message: err.Error(),
			},
		}, nil

	case errors.Is(err, application.ErrActiveAlreadyExists):
		return oapi.CreateApplication409JSONResponse{
			Error:   "conflict",
			Message: err.Error(),
		}, nil

	default:
		return nil, err
	}
}

func handleCandViewErr(err error) (oapi.GetMyApplicationResponseObject, error) {
	switch {
	case errors.Is(err, application.ErrResumeAccessDenied),
		errors.Is(err, identity.ErrCandidateRoleRequired):
		return oapi.GetMyApplication403JSONResponse{
			ForbiddenErrorJSONResponse: oapi.ForbiddenErrorJSONResponse{
				Error:   "forbidden",
				Message: err.Error(),
			},
		}, nil

	case errors.Is(err, application.ErrNotFound):
		return oapi.GetMyApplication404JSONResponse{
			Error:   "not_found",
			Message: err.Error(),
		}, nil

	default:
		return nil, err
	}
}

func handleListMyAppErr(err error) (oapi.ListMyApplicationsResponseObject, error) {
	switch {
	case errors.Is(err, cursors.ErrUnsupportedOrder),
		errors.Is(err, cursors.ErrInvalidCursor),
		errors.Is(err, cursors.ErrCursorOrderMismatch):
		return oapi.ListMyApplications400JSONResponse{
			BadRequestJSONResponse: oapi.BadRequestJSONResponse{
				Error:   "bad_request",
				Message: err.Error(),
			},
		}, nil

	case errors.Is(err, identity.ErrCandidateRoleRequired):
		return oapi.ListMyApplications403JSONResponse{
			ForbiddenErrorJSONResponse: oapi.ForbiddenErrorJSONResponse{
				Error:   "forbidden",
				Message: err.Error(),
			},
		}, nil

	default:
		return nil, err
	}
}

func handleCandiHistoryErr(err error) (oapi.GetMyApplicationHistoryResponseObject, error) {
	switch {
	case errors.Is(err, identity.ErrCandidateRoleRequired):
		return oapi.GetMyApplicationHistory403JSONResponse{
			ForbiddenErrorJSONResponse: oapi.ForbiddenErrorJSONResponse{
				Error:   "forbidden",
				Message: err.Error(),
			},
		}, nil

	case errors.Is(err, application.ErrNotFound):
		return oapi.GetMyApplicationHistory404JSONResponse{
			Error:   "not_found",
			Message: err.Error(),
		}, nil
	}

	return nil, err
}

func handleWithdrawErr(err error) (oapi.WithdrawApplicationResponseObject, error) {
	switch {
	case errors.Is(err, application.ErrCoverLetterTooLong):
		return oapi.WithdrawApplication400JSONResponse{
			BadRequestJSONResponse: oapi.BadRequestJSONResponse{
				Error:   "bad_request",
				Message: err.Error(),
			},
		}, nil

	case errors.Is(err, identity.ErrCandidateRoleRequired):
		return oapi.WithdrawApplication403JSONResponse{
			ForbiddenErrorJSONResponse: oapi.ForbiddenErrorJSONResponse{
				Error:   "forbidden",
				Message: err.Error(),
			},
		}, nil

	case errors.Is(err, application.ErrNotFound):
		return oapi.WithdrawApplication404JSONResponse{
			Error:   "not_found",
			Message: err.Error(),
		}, nil

	case errors.Is(err, application.ErrInvalidStatusTransition),
		errors.Is(err, application.ErrStatusAlreadySet):
		return oapi.WithdrawApplication409JSONResponse{
			Error:   "conflict",
			Message: err.Error(),
		}, nil
	}

	return nil, err
}

func handleHrCompListErr(err error) (oapi.ListCompanyApplicationsResponseObject, error) {
	switch {
	case errors.Is(err, cursors.ErrUnsupportedOrder),
		errors.Is(err, cursors.ErrInvalidCursor),
		errors.Is(err, cursors.ErrCursorOrderMismatch),
		errors.Is(err, listhrsummary.ErrCompanyOrVacancyRequired):
		return oapi.ListCompanyApplications400JSONResponse{
			BadRequestJSONResponse: oapi.BadRequestJSONResponse{
				Error:   "bad_request",
				Message: err.Error(),
			},
		}, nil

	case errors.Is(err, application.ErrAccessDenied),
		errors.Is(err, identity.ErrHrRoleRequired):
		return oapi.ListCompanyApplications403JSONResponse{
			ForbiddenErrorJSONResponse: oapi.ForbiddenErrorJSONResponse{
				Error:   "forbidden",
				Message: err.Error(),
			},
		}, nil

	case errors.Is(err, projection.ErrCompanyNotFound),
		errors.Is(err, projection.ErrVacancyNotFound):
		return oapi.ListCompanyApplications404JSONResponse{
			Error:   "not_found",
			Message: err.Error(),
		}, nil

	default:
		return nil, err
	}
}

func handleHrVacListErr(err error) (oapi.ListVacancyApplicationsResponseObject, error) {
	switch {
	case errors.Is(err, cursors.ErrUnsupportedOrder),
		errors.Is(err, cursors.ErrInvalidCursor),
		errors.Is(err, cursors.ErrCursorOrderMismatch),
		errors.Is(err, listhrsummary.ErrCompanyOrVacancyRequired):
		return oapi.ListVacancyApplications400JSONResponse{
			BadRequestJSONResponse: oapi.BadRequestJSONResponse{
				Error:   "bad_request",
				Message: err.Error(),
			},
		}, nil

	case errors.Is(err, application.ErrAccessDenied),
		errors.Is(err, identity.ErrHrRoleRequired):
		return oapi.ListVacancyApplications403JSONResponse{
			ForbiddenErrorJSONResponse: oapi.ForbiddenErrorJSONResponse{
				Error:   "forbidden",
				Message: err.Error(),
			},
		}, nil

	case errors.Is(err, projection.ErrVacancyNotFound):
		return oapi.ListVacancyApplications404JSONResponse{
			Error:   "not_found",
			Message: err.Error(),
		}, nil

	default:
		return nil, err
	}
}

func handleHrViewErr(err error) (oapi.GetHrApplicationResponseObject, error) {
	switch {
	case errors.Is(err, identity.ErrHrRoleRequired):
		return oapi.GetHrApplication403JSONResponse{
			ForbiddenErrorJSONResponse: oapi.ForbiddenErrorJSONResponse{
				Error:   "forbidden",
				Message: err.Error(),
			},
		}, nil

	case errors.Is(err, application.ErrNotFound):
		return oapi.GetHrApplication404JSONResponse{
			Error:   "not_found",
			Message: err.Error(),
		}, nil
	}

	return nil, err
}

func handleHrUpdStatus(err error) (oapi.ChangeApplicationStatusResponseObject, error) {
	switch {
	case errors.Is(err, application.ErrInvalidStatus),
		errors.Is(err, application.ErrCoverLetterTooLong):
		return oapi.ChangeApplicationStatus400JSONResponse{
			BadRequestJSONResponse: oapi.BadRequestJSONResponse{
				Error:   "bad_request",
				Message: err.Error(),
			},
		}, nil

	case errors.Is(err, identity.ErrHrRoleRequired):
		return oapi.ChangeApplicationStatus403JSONResponse{
			ForbiddenErrorJSONResponse: oapi.ForbiddenErrorJSONResponse{
				Error:   "forbidden",
				Message: err.Error(),
			},
		}, nil

	case errors.Is(err, application.ErrNotFound):
		return oapi.ChangeApplicationStatus404JSONResponse{
			Error:   "not_found",
			Message: err.Error(),
		}, nil

	case errors.Is(err, application.ErrInvalidStatusTransition),
		errors.Is(err, application.ErrStatusAlreadySet):
		return oapi.ChangeApplicationStatus409JSONResponse{
			Error:   "conflict",
			Message: err.Error(),
		}, nil
	}

	return nil, err
}

func handleGetHistoryHrViewErr(err error) (oapi.GetHrApplicationHistoryResponseObject, error) {
	switch {
	case errors.Is(err, identity.ErrHrRoleRequired):
		return oapi.GetHrApplicationHistory403JSONResponse{
			ForbiddenErrorJSONResponse: oapi.ForbiddenErrorJSONResponse{
				Error:   "forbidden",
				Message: err.Error(),
			},
		}, nil

	case errors.Is(err, application.ErrNotFound):
		return oapi.GetHrApplicationHistory404JSONResponse{
			Error:   "not_found",
			Message: err.Error(),
		}, nil
	}

	return nil, err
}

func handleCompanyAnalyticsSummary(err error) (oapi.GetCompanyAnalyticsSummaryResponseObject, error) {
	switch {
	case errors.Is(err, identity.ErrHrRoleRequired):
		return oapi.GetCompanyAnalyticsSummary403JSONResponse{
			ForbiddenErrorJSONResponse: oapi.ForbiddenErrorJSONResponse{
				Error:   "forbidden",
				Message: err.Error(),
			},
		}, nil

	case errors.Is(err, projection.ErrCompanyNotFound):
		return oapi.GetCompanyAnalyticsSummary404JSONResponse{
			Error:   "not_found",
			Message: err.Error(),
		}, nil
	}

	return nil, err
}

func handleVacancyAnalyticsSummaryErr(err error) (oapi.GetVacancyAnalyticsSummaryResponseObject, error) {
	switch {
	case errors.Is(err, identity.ErrHrRoleRequired):
		return oapi.GetVacancyAnalyticsSummary403JSONResponse{
			ForbiddenErrorJSONResponse: oapi.ForbiddenErrorJSONResponse{
				Error:   "forbidden",
				Message: err.Error(),
			},
		}, nil

	case errors.Is(err, projection.ErrVacancyNotFound):
		return oapi.GetVacancyAnalyticsSummary404JSONResponse{
			Error:   "not_found",
			Message: err.Error(),
		}, nil
	}

	return nil, err
}

func handleCompanyDynamicsErr(err error) (oapi.GetCompanyDynamicsResponseObject, error) {
	switch {
	case errors.Is(err, dynamics.ErrInvalidPeriod),
		errors.Is(err, dynamics.ErrPeriodTooLarge),
		errors.Is(err, dynamics.ErrInvalidInterval):
		return oapi.GetCompanyDynamics400JSONResponse{
			BadRequestJSONResponse: oapi.BadRequestJSONResponse{
				Error:   "forbidden",
				Message: err.Error(),
			},
		}, nil

	case errors.Is(err, identity.ErrHrRoleRequired):
		return oapi.GetCompanyDynamics403JSONResponse{
			ForbiddenErrorJSONResponse: oapi.ForbiddenErrorJSONResponse{
				Error:   "forbidden",
				Message: err.Error(),
			},
		}, nil

	case errors.Is(err, projection.ErrCompanyNotFound):
		return oapi.GetCompanyDynamics404JSONResponse{
			Error:   "not_found",
			Message: err.Error(),
		}, nil
	}

	return nil, err
}

func handleVacancyDynamicsErr(err error) (oapi.GetVacancyDynamicsResponseObject, error) {
	switch {
	case errors.Is(err, dynamics.ErrInvalidPeriod),
		errors.Is(err, dynamics.ErrPeriodTooLarge),
		errors.Is(err, dynamics.ErrInvalidInterval):
		return oapi.GetVacancyDynamics400JSONResponse{
			BadRequestJSONResponse: oapi.BadRequestJSONResponse{
				Error:   "forbidden",
				Message: err.Error(),
			},
		}, nil

	case errors.Is(err, identity.ErrHrRoleRequired):
		return oapi.GetVacancyDynamics403JSONResponse{
			ForbiddenErrorJSONResponse: oapi.ForbiddenErrorJSONResponse{
				Error:   "forbidden",
				Message: err.Error(),
			},
		}, nil

	case errors.Is(err, projection.ErrVacancyNotFound):
		return oapi.GetVacancyDynamics404JSONResponse{
			Error:   "not_found",
			Message: err.Error(),
		}, nil
	}

	return nil, err
}
