package value_types

type CompanyRole string

const (
	CompanyRoleRecruiter CompanyRole = "recruiter"
	CompanyRoleAdmin     CompanyRole = "admin"
)

func (r CompanyRole) IsValid() bool {
	switch r {
	case CompanyRoleRecruiter, CompanyRoleAdmin:
		return true
	}

	return false
}
