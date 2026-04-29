package publish

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

type outboxWriter interface {
	WriteVacancyPublished(ctx context.Context, ev vacancy.PublishedEvent) error
}
