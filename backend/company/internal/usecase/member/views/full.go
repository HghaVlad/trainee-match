package views

import "github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"

type MemberFullView struct {
	Member   member.CompanyMember
	Username string
	Email    string
}
