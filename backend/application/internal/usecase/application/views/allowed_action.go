package views

import "github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"

type AllowedAction string

const (
	AllowedActionWithdraw        AllowedAction = "withdraw"
	AllowedActionMarkSeen        AllowedAction = "mark_seen"
	AllowedActionInviteInterview AllowedAction = "invite_interview"
	AllowedActionReject          AllowedAction = "reject"
	AllowedActionMakeOffer       AllowedAction = "make_offer"
)

func HRAllowedActions(status application.Status) []AllowedAction {
	transitions := status.AvailableTransitions(application.ActorHR)
	res := make([]AllowedAction, 0, len(transitions))

	for _, t := range transitions {
		//nolint:exhaustive // others won't contribute to res
		switch t {
		case application.StatusSeen:
			res = append(res, AllowedActionMarkSeen)

		case application.StatusInterview:
			res = append(res, AllowedActionInviteInterview)

		case application.StatusRejected:
			res = append(res, AllowedActionReject)

		case application.StatusOffer:
			res = append(res, AllowedActionMakeOffer)
		default:
		}
	}

	return res
}

func CandidateAllowedActions(status application.Status) []AllowedAction {
	transitions := status.AvailableTransitions(application.ActorCandidate)
	res := make([]AllowedAction, 0, len(transitions))

	for _, t := range transitions {
		if t == application.StatusWithdrawn {
			res = append(res, AllowedActionWithdraw)
			break
		}
	}

	return res
}
