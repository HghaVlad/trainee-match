package remove

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/company"
)

//go:generate mockgen -source=port.go -destination=mocks/port_mocks.go -package=mocks
type outboxWriter interface {
	WriteCompanyDeleted(ctx context.Context, ev company.DeletedEvent) error
}
