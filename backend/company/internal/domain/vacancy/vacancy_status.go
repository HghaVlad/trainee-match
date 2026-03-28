package vacancy

type Status string

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
	StatusArchived  Status = "archived"
)

func (vs Status) IsValid() bool {
	switch vs {
	case StatusDraft,
		StatusPublished,
		StatusArchived:
		return true
	}

	return false
}
