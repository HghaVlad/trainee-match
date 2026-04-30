package update

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
)

//go:generate mockgen -source=port.go -destination=mocks/port_mocks.go -package=mocks
type outboxWriter interface {
	WriteVacancyUpdated(ctx context.Context, ev vacancy.UpdatedEvent) error
}
