package mappers

import (
	openapitypes "github.com/oapi-codegen/runtime/types"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/oapi"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/apply"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/listcandidatesummary"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
)

func ApplyReqToUC(dto oapi.CreateApplicationRequestObject) apply.Request {
	return apply.Request{
		VacancyID:   dto.Body.VacancyId,
		ResumeID:    dto.Body.ResumeId,
		CoverLetter: dto.Body.CoverLetter,
	}
}

func ListMyApplicationsReqToUC(dto oapi.ListMyApplicationsRequestObject) listcandidatesummary.Request {
	req := listcandidatesummary.Request{
		Limit: 20,
		Order: listcandidatesummary.OrderCreatedAtDesc,
	}

	if dto.Params.Statuses != nil {
		req.Statuses = make([]application.Status, 0, len(*dto.Params.Statuses))
		for _, status := range *dto.Params.Statuses {
			req.Statuses = append(req.Statuses, application.Status(status))
		}
	}

	if dto.Params.CompanyId != nil {
		companyID := *dto.Params.CompanyId
		req.CompanyID = &companyID
	}

	if dto.Params.Cursor != nil {
		req.Cursor = *dto.Params.Cursor
	}

	if dto.Params.Limit != nil {
		req.Limit = *dto.Params.Limit
	}

	if dto.Params.Sort != nil {
		req.Order = listcandidatesummary.Order(*dto.Params.Sort)
	}

	return req
}

func CandidateListResponseToHTTP(resp *listcandidatesummary.Response) oapi.ListMyApplications200JSONResponse {
	items := make([]oapi.CandidateApplicationListItem, 0, len(resp.AppSummaries))

	for _, item := range resp.AppSummaries {
		items = append(items, oapi.CandidateApplicationListItem{
			CompanyId:    item.CompanyID,
			CompanyName:  item.CompanyName,
			CreatedAt:    item.CreatedAt,
			Id:           item.AppID,
			Status:       oapi.ApplicationStatus(item.Status),
			UpdatedAt:    item.UpdatedAt,
			VacancyId:    item.VacancyID,
			VacancyTitle: item.VacancyTitle,
		})
	}

	return oapi.ListMyApplications200JSONResponse{
		Data:       items,
		HasNext:    resp.HasNext,
		NextCursor: resp.NextCursor,
	}
}

func CandidateViewWithDetailsToHTTP(view *views.CandidateViewWithDetails) oapi.CandidateApplicationDetails {
	return oapi.CandidateApplicationDetails{
		Id:             view.AppID,
		Status:         oapi.ApplicationStatus(view.Status),
		VacancyId:      view.VacancyID,
		VacancyTitle:   view.VacancyTitle,
		CompanyId:      view.CompanyID,
		CompanyName:    view.CompanyName,
		CoverLetter:    view.CoverLetter,
		AllowedActions: mapCandidateAllowedActions(view.AllowedActions),
		CreatedAt:      view.CreatedAt,
		UpdatedAt:      view.UpdatedAt,

		Snapshot: oapi.ApplicationSnapshot{
			Email:     emailToOAPI(&view.Snapshot.Email),
			FullName:  view.Snapshot.FullName,
			Telegram:  view.Snapshot.Telegram,
			CreatedAt: view.Snapshot.CreatedAt,
			ResumeData: map[string]any{
				"resumeData": view.Snapshot.ResumeData,
			},
		},

		StatusHistory: mapCandidateHistory(view.StatusHistory),
	}
}

func mapCandidateAllowedActions(actions []views.AllowedAction) []oapi.CandidateAllowedAction {
	result := make([]oapi.CandidateAllowedAction, 0, len(actions))

	for _, action := range actions {
		result = append(result, oapi.CandidateAllowedAction(action))
	}

	return result
}

func mapCandidateHistory(items []views.StatusChangeCandidateView) []oapi.CandidateApplicationStatusHistoryItem {
	result := make([]oapi.CandidateApplicationStatusHistoryItem, 0, len(items))

	for _, item := range items {
		result = append(result, oapi.CandidateApplicationStatusHistoryItem{
			ChangedByRole: oapi.CandidateApplicationStatusHistoryItemChangedByRole(item.ChangedByRole),
			CreatedAt:     item.CreatedAt,
			Status:        oapi.ApplicationStatus(item.Status),
		})
	}

	return result
}

func emailToOAPI(email *string) *openapitypes.Email {
	if email == nil {
		return nil
	}

	typed := openapitypes.Email(*email)
	return &typed
}
