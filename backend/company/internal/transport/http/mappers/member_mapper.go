package mappers

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/dto"
	addmemb "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/add"
	update_member "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/member/update"
)

func CompanyAddHrReqToUC(
	companyID uuid.UUID,
	dtoReq *dto.CompanyAddHrRequest,
) *addmemb.Request {
	return &addmemb.Request{
		CompanyID: companyID,
		UserID:    dtoReq.UserID,
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
