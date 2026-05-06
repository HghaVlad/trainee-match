package add

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/member"
)

//go:generate mockgen -source=port.go -destination=mocks/port_mocks.go -package=mocks
type outboxWriter interface {
	WriteCompanyMemberAdded(ctx context.Context, ev member.AddedEvent) error
}
