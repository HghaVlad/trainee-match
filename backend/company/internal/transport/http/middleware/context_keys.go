package middleware

import (
	"context"

	uc_common "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common/identity"
)

type ctxIdentityKeyT struct{}

//nolint:gochecknoglobals // ctx key
var identityKey = ctxIdentityKeyT{}

func WithIdentity(ctx context.Context, identity uc_common.Identity) context.Context {
	return context.WithValue(ctx, identityKey, identity)
}

func IdentityFromContext(ctx context.Context) uc_common.Identity {
	id, _ := ctx.Value(identityKey).(uc_common.Identity)
	return id
}
