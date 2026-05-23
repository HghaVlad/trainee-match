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

type ModerationStatus string

const (
	ModerationStatusOK     ModerationStatus = "ok"
	ModerationStatusHidden ModerationStatus = "hidden"
)

func (ms ModerationStatus) IsValid() bool {
	switch ms {
	case ModerationStatusOK,
		ModerationStatusHidden:
		return true
	}

	return false
}
