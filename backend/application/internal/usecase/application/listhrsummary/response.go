package listhrsummary

import "github.com/HghaVlad/trainee-match/backend/application/internal/usecase/application/views"

type Response struct {
	AppSummaries []views.HrAppSummary
	NextCursor   *string
	HasNext      bool
}
