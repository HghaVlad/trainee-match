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
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/identity"
)

type Handler struct {
	apply             applyUC
	listCandidateApps listCandidateAppsUC
	getCandidateView  getCandidateViewUC
	logger            *slog.Logger
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		apply:             deps.Apply,
		listCandidateApps: deps.ListCandidateApps,
		getCandidateView:  deps.GetCandidateViewUC,
		logger:            deps.Logger,
	}
}

func (h *Handler) ListMyApplications(
	ctx context.Context,
	request oapi.ListMyApplicationsRequestObject,
) (oapi.ListMyApplicationsResponseObject, error) {
	ident := middleware.IdentityFromContext(ctx)
	req := mappers.ListMyApplicationsReqToUC(request)

	resp, err := h.listCandidateApps.Execute(ctx, req, *ident)
	if err != nil {
		return nil, err
	}

	return mappers.CandidateListResponseToHTTP(resp), nil
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
		Data: mappers.CandidateViewWithDetailsToHTTP(resp),
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
		return getCandViewErrToResponse(err)
	}

	return oapi.GetMyApplication200JSONResponse{
		Data: mappers.CandidateViewWithDetailsToHTTP(resp),
	}, nil
}

func (h *Handler) GetMyApplicationHistory(
	ctx context.Context,
	request oapi.GetMyApplicationHistoryRequestObject,
) (oapi.GetMyApplicationHistoryResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) WithdrawApplication(
	ctx context.Context,
	request oapi.WithdrawApplicationRequestObject,
) (oapi.WithdrawApplicationResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) GetHrApplication(
	ctx context.Context,
	request oapi.GetHrApplicationRequestObject,
) (oapi.GetHrApplicationResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) GetHrApplicationHistory(
	ctx context.Context,
	request oapi.GetHrApplicationHistoryRequestObject,
) (oapi.GetHrApplicationHistoryResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) ChangeApplicationStatus(
	ctx context.Context,
	request oapi.ChangeApplicationStatusRequestObject,
) (oapi.ChangeApplicationStatusResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) GetCompanyDynamics(
	ctx context.Context,
	request oapi.GetCompanyDynamicsRequestObject,
) (oapi.GetCompanyDynamicsResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) GetCompanyStatusFunnel(
	ctx context.Context,
	request oapi.GetCompanyStatusFunnelRequestObject,
) (oapi.GetCompanyStatusFunnelResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) GetCompanyAnalyticsSummary(
	ctx context.Context,
	request oapi.GetCompanyAnalyticsSummaryRequestObject,
) (oapi.GetCompanyAnalyticsSummaryResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) ListCompanyApplications(
	ctx context.Context,
	request oapi.ListCompanyApplicationsRequestObject,
) (oapi.ListCompanyApplicationsResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) GetVacancyDynamics(
	ctx context.Context,
	request oapi.GetVacancyDynamicsRequestObject,
) (oapi.GetVacancyDynamicsResponseObject, error) {
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

func (h *Handler) GetVacancyAnalyticsSummary(
	ctx context.Context,
	request oapi.GetVacancyAnalyticsSummaryRequestObject,
) (oapi.GetVacancyAnalyticsSummaryResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) ListVacancyApplications(
	ctx context.Context,
	request oapi.ListVacancyApplicationsRequestObject,
) (oapi.ListVacancyApplicationsResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func applyErrToResponse(err error) (oapi.CreateApplicationResponseObject, error) {
	switch {
	case errors.Is(err, application.ErrResumeAccessDenied),
		errors.Is(err, identity.ErrCandidateRoleRequired):
		return oapi.CreateApplication403Response{}, nil

	case errors.Is(err, projection.ErrVacancyNotFound),
		errors.Is(err, projection.ErrResumeNotFound),
		errors.Is(err, projection.ErrCandidateNotFound):
		return oapi.CreateApplication404Response{}, nil

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

func getCandViewErrToResponse(err error) (oapi.GetMyApplicationResponseObject, error) {
	switch {
	case errors.Is(err, application.ErrResumeAccessDenied),
		errors.Is(err, identity.ErrCandidateRoleRequired):
		return oapi.GetMyApplication403Response{}, nil

	case errors.Is(err, application.ErrNotFound):
		return oapi.GetMyApplication404Response{}, nil

	default:
		return nil, err
	}
}
