package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/helpers"
)

type ctxBodyKeyT struct{}

//nolint:gochecknoglobals // ctx key
var ctxBodyKey = ctxBodyKeyT{}

// BindJSONBodyMiddleware binds json body to the given type T.
// Responds with responds.ErrorResponse 400 if binding failed.
// Saves value to request context, use BodyFromContext to get it
func BindJSONBodyMiddleware[T any]() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				_ = r.Body.Close()
			}()

			var body T
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				helpers.RespondErrorMsg(w, http.StatusBadRequest, "invalid JSON body")
				return
			}

			ctx := context.WithValue(r.Context(), ctxBodyKey, &body)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// TODO: test behaviour with pointers

// BodyFromContext retrieves typed object from the context.
// It implies that BodyFromContext was used
func BodyFromContext[T any](ctx context.Context) *T {
	body, ok := ctx.Value(ctxBodyKey).(*T)
	if !ok {
		panic("body not found in context: bind json body middleware is not applied")
	}

	return body
}
