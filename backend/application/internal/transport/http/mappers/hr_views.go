package mappers

import (
	"github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/oapi"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"
)

func HrDetailedViewToHTTP(view *views.HrDetailedView) oapi.HrApplicationDetails {
	return oapi.HrApplicationDetails{
		Id:             view.AppID,
		Status:         oapi.ApplicationStatus(view.Status),
		VacancyId:      view.VacancyID,
		VacancyTitle:   view.VacancyTitle,
		CoverLetter:    view.CoverLetter,
		AllowedActions: mapHrAllowedActions(view.AllowedActions),
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

		StatusHistory: HrHistoryToHTTP(view.StatusHistory),
	}
}

func HrHistoryToHTTP(items []views.StatusChangeHrFullView) []oapi.HrApplicationStatusHistoryItem {
	result := make([]oapi.HrApplicationStatusHistoryItem, 0, len(items))

	for _, item := range items {
		result = append(result, oapi.HrApplicationStatusHistoryItem{
			ChangedByRole:   oapi.ActorType(item.ChangedByRole),
			ChangedByUserId: item.ChangedByUserID,
			Comment:         item.Comment,
			CreatedAt:       item.CreatedAt,
			Status:          oapi.ApplicationStatus(item.Status),
		})
	}

	return result
}

func mapHrAllowedActions(actions []views.AllowedAction) []oapi.HrAllowedAction {
	result := make([]oapi.HrAllowedAction, 0, len(actions))

	for _, action := range actions {
		result = append(result, oapi.HrAllowedAction(action))
	}

	return result
}
