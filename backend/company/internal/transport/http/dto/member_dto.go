package dto

import (
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
)

type CompanyAddHrRequest struct {
	UserID uuid.UUID `json:"userId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Role   string    `json:"role"   example:"recruiter"                            enums:"recruiter,admin"`
}

type CompanyUpdateMemberRequest struct {
	Role string `json:"role" enums:"recruiter,admin" example:"admin"`
}

type CompanyMemberListItem struct {
	UserID    uuid.UUID          `json:"userId"    example:"550e8400-e29b-41d4-a716-446655440000"`
	CompanyID uuid.UUID          `json:"companyId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username  string             `json:"username"  example:"JohnKaisen"`
	Email     string             `json:"email"     example:"johnkaisen@gmail.com"`
	Role      member.CompanyRole `json:"role"      example:"recruiter"                            enums:"recruiter,admin"`
}

type CompanyMemberListResponse struct {
	Members []CompanyMemberListItem `json:"members"`
	HasMore bool                    `json:"hasMore"`
}
