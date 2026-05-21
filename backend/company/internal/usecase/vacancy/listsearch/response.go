package listsearch

import (
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/views"
)

type Response struct {
	Vacancies  []views.PublishedVacSummary
	NextCursor *string
}
