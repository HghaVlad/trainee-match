package handler

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/oapi"
)

type Handler struct {
}

func (h *Handler) ListMyApplications(
	ctx context.Context,
	request oapi.ListMyApplicationsRequestObject,
) (oapi.ListMyApplicationsResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) CreateApplication(
	ctx context.Context,
	request oapi.CreateApplicationRequestObject,
) (oapi.CreateApplicationResponseObject, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) GetMyApplication(
	ctx context.Context,
	request oapi.GetMyApplicationRequestObject,
) (oapi.GetMyApplicationResponseObject, error) {
	// TODO implement me
	panic("implement me")
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
