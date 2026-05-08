package application

type Status string

const (
	StatusSubmitted Status = "submitted"
	StatusSeen      Status = "seen"
	StatusInterview Status = "interview"
	StatusRejected  Status = "rejected"
	StatusOffer     Status = "offer"
	StatusWithdrawn Status = "withdrawn"
)

//nolint:gochecknoglobals // its readonly and for validations
var allowedTransitions = map[Status]map[Status]struct{}{
	StatusSubmitted: {
		StatusSeen:      {},
		StatusRejected:  {},
		StatusWithdrawn: {},
	},

	StatusSeen: {
		StatusInterview: {},
		StatusRejected:  {},
		StatusWithdrawn: {},
	},

	StatusInterview: {
		StatusOffer:     {},
		StatusRejected:  {},
		StatusWithdrawn: {},
	},

	StatusOffer: {},

	StatusRejected: {},

	StatusWithdrawn: {},
}

func (s Status) CanTransitionTo(next Status) bool {
	allowed, ok := allowedTransitions[s]
	if !ok {
		return false
	}

	_, ok = allowed[next]

	return ok
}
