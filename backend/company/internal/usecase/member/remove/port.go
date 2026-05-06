package remove

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
)

//go:generate mockgen -source=port.go -destination=mocks/port_mocks.go -package=mocks
type outboxWriter interface {
	WriteCompanyMemberRemoved(ctx context.Context, ev member.RemovedEvent) error
}
