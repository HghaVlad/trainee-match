package mappers

import (
	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
	"github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/oapi"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/listcandidatesummary"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/listhrsummary"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/common/cursors"
)

func ListMyApplicationsReqToUC(dto oapi.ListMyApplicationsRequestObject) listcandidatesummary.Request {
	req := listcandidatesummary.Request{
		Limit: 20,
		Order: cursors.OrderCreatedAtDesc,
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
		req.Order = cursors.SummaryOrder(*dto.Params.Sort)
	}

	return req
}

func ListCompanyApplicationsReqToUC(dto oapi.ListCompanyApplicationsRequestObject) listhrsummary.Request {
	req := listhrsummary.Request{
		Limit: 20,
		Order: cursors.HrSummaryOrderCreatedAtDesc,
	}

	companyID := dto.CompanyId
	req.CompanyID = &companyID

	if dto.Params.Statuses != nil {
		req.Statuses = make([]application.Status, 0, len(*dto.Params.Statuses))
		for _, status := range *dto.Params.Statuses {
			req.Statuses = append(req.Statuses, application.Status(status))
		}
	}

	if dto.Params.VacancyId != nil {
		req.VacancyID = dto.Params.VacancyId
	}

	if dto.Params.CreatedFrom != nil {
		req.CreatedFrom = dto.Params.CreatedFrom
	}

	if dto.Params.CreatedTo != nil {
		req.CreatedTo = dto.Params.CreatedTo
	}

	if dto.Params.Cursor != nil {
		req.Cursor = *dto.Params.Cursor
	}

	if dto.Params.Limit != nil {
		req.Limit = *dto.Params.Limit
	}

	if dto.Params.Sort != nil {
		req.Order = cursors.HrSummaryOrder(*dto.Params.Sort)
	}

	return req
}

func ListVacancyApplicationsReqToUC(dto oapi.ListVacancyApplicationsRequestObject) listhrsummary.Request {
	req := listhrsummary.Request{
		Limit: 20,
		Order: cursors.HrSummaryOrderCreatedAtDesc,
	}

	req.VacancyID = &dto.VacancyId

	if dto.Params.Statuses != nil {
		req.Statuses = make([]application.Status, 0, len(*dto.Params.Statuses))
		for _, status := range *dto.Params.Statuses {
			req.Statuses = append(req.Statuses, application.Status(status))
		}
	}

	if dto.Params.CreatedFrom != nil {
		req.CreatedFrom = dto.Params.CreatedFrom
	}

	if dto.Params.CreatedTo != nil {
		req.CreatedTo = dto.Params.CreatedTo
	}

	if dto.Params.Cursor != nil {
		req.Cursor = *dto.Params.Cursor
	}

	if dto.Params.Limit != nil {
		req.Limit = *dto.Params.Limit
	}

	if dto.Params.Sort != nil {
		req.Order = cursors.HrSummaryOrder(*dto.Params.Sort)
	}

	return req
}

func HrListResponseToHTTP(resp *listhrsummary.Response) oapi.HrApplicationListResponse {
	items := make([]oapi.HrApplicationListItem, 0, len(resp.AppSummaries))

	for _, item := range resp.AppSummaries {
		items = append(items, oapi.HrApplicationListItem{
			CreatedAt: item.CreatedAt,
			Id:        item.AppID,
			Snapshot: oapi.ApplicationSnapshotSummary{
				CreatedAt: item.AppSnap.CreatedAt,
				Email:     emailToOAPI(&item.AppSnap.Email),
				FullName:  item.AppSnap.FullName,
				Telegram:  item.AppSnap.Telegram,
			},
			Status:       oapi.ApplicationStatus(item.Status),
			UpdatedAt:    item.UpdatedAt,
			VacancyId:    item.VacancyID,
			VacancyTitle: item.VacancyTitle,
		})
	}

	return oapi.HrApplicationListResponse{
		Data:       items,
		HasNext:    resp.HasNext,
		NextCursor: resp.NextCursor,
	}
}
