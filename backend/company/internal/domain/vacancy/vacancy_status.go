package vacancy

type VacancyStatus string

const (
	VacancyStatusDraft     VacancyStatus = "draft"
	VacancyStatusPublished VacancyStatus = "published"
	VacancyStatusArchived  VacancyStatus = "archived"
)

func (vs VacancyStatus) IsValid() bool {
	switch vs {
	case VacancyStatusDraft,
		VacancyStatusPublished,
		VacancyStatusArchived:
		return true
	}

	return false
}
