package value_types

type WorkFormat string

const (
	WorkFormatOnSite WorkFormat = "onsite"
	WorkFormatRemote WorkFormat = "remote"
	WorkFormatHybrid WorkFormat = "hybrid"
)

func (wf WorkFormat) IsValid() bool {
	switch wf {
	case WorkFormatOnSite,
		WorkFormatRemote,
		WorkFormatHybrid:
		return true
	default:
		return false
	}
}
