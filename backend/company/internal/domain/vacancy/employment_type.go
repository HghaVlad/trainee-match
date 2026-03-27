package vacancy

type EmploymentType string

const (
	EmploymentTypeInternship EmploymentType = "internship"
	EmploymentTypeFullTime   EmploymentType = "full_time"
	EmploymentTypePartTime   EmploymentType = "part_time"
)

func (et EmploymentType) IsValid() bool {
	switch et {
	case EmploymentTypeInternship,
		EmploymentTypeFullTime,
		EmploymentTypePartTime:
		return true
	default:
		return false
	}
}
