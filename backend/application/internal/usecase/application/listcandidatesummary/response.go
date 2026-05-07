package listcandidatesummary

import "github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"

type Response struct {
	AppSummaries []views.CandidateAppSummary
	NextCursor   *string
	HasNext      bool
}
