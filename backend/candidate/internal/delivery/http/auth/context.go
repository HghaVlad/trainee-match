package auth

import "context"

type contextKey string

func WithUser(ctx context.Context, u User) context.Context {
	return context.WithValue(ctx, "user", u)
}

func FromContext(ctx context.Context) (User, bool) {
	u, ok := ctx.Value("user").(User)
	return u, ok
}
