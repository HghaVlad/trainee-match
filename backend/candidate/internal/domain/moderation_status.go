package domain

type ModerationStatus string

const (
	ModerationStatusOK     ModerationStatus = "ok"
	ModerationStatusHidden ModerationStatus = "hidden"
)

func (s ModerationStatus) IsValid() bool {
	switch s {
	case ModerationStatusOK, ModerationStatusHidden:
		return true
	default:
		return false
	}
}
