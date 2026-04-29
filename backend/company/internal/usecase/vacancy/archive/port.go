package archive

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

type outboxWriter interface {
	WriteVacancyArchived(ctx context.Context, ev vacancy.ArchivedEvent) error
}
