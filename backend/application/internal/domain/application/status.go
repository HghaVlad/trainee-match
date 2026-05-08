package application

import "slices"

type Status string

const (
	StatusSubmitted Status = "submitted"
	StatusSeen      Status = "seen"
	StatusInterview Status = "interview"
	StatusRejected  Status = "rejected"
	StatusOffer     Status = "offer"
	StatusWithdrawn Status = "withdrawn"
)

type Transition struct {
	To     Status
	Actors []Actor
}

//nolint:gochecknoglobals // its readonly and for validations
var transitions = map[Status][]Transition{
	StatusSubmitted: {
		{To: StatusSeen, Actors: []Actor{ActorHR}},
		{To: StatusInterview, Actors: []Actor{ActorHR}},
		{To: StatusRejected, Actors: []Actor{ActorHR}},
		{To: StatusWithdrawn, Actors: []Actor{ActorCandidate}},
	},

	StatusSeen: {
		{To: StatusInterview, Actors: []Actor{ActorHR}},
		{To: StatusRejected, Actors: []Actor{ActorHR}},
		{To: StatusWithdrawn, Actors: []Actor{ActorCandidate}},
	},

	StatusInterview: {
		{To: StatusOffer, Actors: []Actor{ActorHR}},
		{To: StatusRejected, Actors: []Actor{ActorHR}},
		{To: StatusWithdrawn, Actors: []Actor{ActorCandidate}},
	},

	StatusOffer:     {},
	StatusRejected:  {},
	StatusWithdrawn: {},
}

func (s Status) CanTransitionTo(next Status, actor Actor) bool {
	available, ok := transitions[s]
	if !ok {
		return false
	}

	for _, t := range available {
		if t.To != next {
			continue
		}

		if slices.Contains(t.Actors, actor) {
			return true
		}
	}

	return false
}

func (s Status) AvailableTransitions(actor Actor) []Status {
	available, ok := transitions[s]
	if !ok {
		return nil
	}

	res := make([]Status, 0, len(available))

	for _, t := range available {
		if slices.Contains(t.Actors, actor) {
			res = append(res, t.To)
		}
	}

	return res
}
