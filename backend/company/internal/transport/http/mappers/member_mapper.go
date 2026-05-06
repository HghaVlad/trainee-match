package mappers

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/dto"
	addmemb "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/add"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/list"
	update_member "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/update"
)

func CompanyAddHrReqToUC(
	companyID uuid.UUID,
	dtoReq *dto.CompanyAddHrRequest,
) *addmemb.Request {
	return &addmemb.Request{
		CompanyID: companyID,
		Username:  dtoReq.Username,
		Role:      member.CompanyRole(dtoReq.Role),
	}
}

func CompanyUpdateMemberReqToUC(
	companyID uuid.UUID,
	userID uuid.UUID,
	dtoReq *dto.CompanyUpdateMemberRequest,
) *update_member.Request {
	return &update_member.Request{
		CompanyID: companyID,
		UserID:    userID,
		Role:      member.CompanyRole(dtoReq.Role),
	}
}

func CompanyMemberListRespToDto(resp list.Response) dto.CompanyMemberListResponse {
	members := make([]dto.CompanyMemberListItem, len(resp.Members))

	for i, m := range resp.Members {
		members[i] = dto.CompanyMemberListItem{
			UserID:    m.UserID,
			CompanyID: m.CompanyID,
			Username:  m.Username,
			Email:     m.Email,
			Role:      m.Role,
		}
	}

	return dto.CompanyMemberListResponse{
		Members: members,
		HasMore: resp.HasMore,
	}
}
