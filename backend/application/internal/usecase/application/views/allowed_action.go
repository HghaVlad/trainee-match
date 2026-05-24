package views

import (
	"slices"

	"github.com/HghaVlad/trainee-match/backend/application/internal/domain/application"
)

type AllowedAction string

const (
	AllowedActionWithdraw        AllowedAction = "withdraw"
	AllowedActionMarkSeen        AllowedAction = "markSeen"
	AllowedActionInviteInterview AllowedAction = "moveToInterview"
	AllowedActionReject          AllowedAction = "reject"
	AllowedActionMakeOffer       AllowedAction = "makeOffer"
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

	if slices.Contains(transitions, application.StatusWithdrawn) {
		res = append(res, AllowedActionWithdraw)
	}

	return res
}
