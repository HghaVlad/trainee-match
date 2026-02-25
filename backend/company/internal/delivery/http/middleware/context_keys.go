package my_middleware

import (
	"context"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
)

type ctxKey struct{}

var (
	identityKey ctxKey = ctxKey{}
)

func WithIdentity(ctx context.Context, identity uc_common.Identity) context.Context {
	return context.WithValue(ctx, identityKey, identity)
}

func IdentityFromContext(ctx context.Context) uc_common.Identity {
	id := ctx.Value(identityKey).(uc_common.Identity)
	return id
}
