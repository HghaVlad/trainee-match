package domain

type User struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Username  string
	Role      string
}

type UserRole string

const (
	UserCandidateRole UserRole = "Candidate"
	UserCompanyRole   UserRole = "Company"
	UserAdminRole     UserRole = "admin"
)
