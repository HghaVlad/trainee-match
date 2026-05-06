package mappers

import (
	openapitypes "github.com/oapi-codegen/runtime/types"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/oapi"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/apply"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
)

func ApplyReqToUC(dto oapi.CreateApplicationRequestObject) apply.Request {
	return apply.Request{
		VacancyID:   dto.Body.VacancyId,
		ResumeID:    dto.Body.ResumeId,
		CoverLetter: dto.Body.CoverLetter,
	}
}

func ApplyRespToHTTP(view *views.Details) oapi.CandidateApplicationDetails {
	return oapi.CandidateApplicationDetails{
		Id:             view.Application.ID,
		Status:         oapi.ApplicationStatus(view.Application.Status),
		VacancyId:      view.Application.VacancyID,
		VacancyTitle:   view.VacProj.Title,
		CompanyId:      view.VacProj.CompanyID,
		CompanyName:    view.VacProj.CompanyName,
		CoverLetter:    view.Application.CoverLetter,
		AllowedActions: mapCandidateAllowedActions(view.AllowedActions),
		CreatedAt:      view.Application.CreatedAt,
		UpdatedAt:      view.Application.UpdatedAt,

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

func mapCandidateHistory(items []application.StatusChange) []oapi.CandidateApplicationStatusHistoryItem {
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
