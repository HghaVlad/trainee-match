package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	utilslog "github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/utils/logger"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/helpers"
)

type uuidCtxKey struct {
	key string
}

// UUIDMiddleware parses uuid from path key.
// Responds with responds.ErrorResponse 400 if uuid parsing failed.
// Saves value to request context, use UUIDFromContext to get it
func UUIDMiddleware(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			raw := chi.URLParam(r, key)

			id, err := uuid.Parse(raw)
			if err != nil {
				logger := utilslog.FromContext(ctx)
				logger.InfoContext(ctx, "Invalid UUID: "+raw)
				helpers.RespondErrorMsg(ctx, w, http.StatusBadRequest, "Invalid UUID: "+raw)
				return
			}

			ctxKey := uuidCtxKey{key: key}
			ctx = context.WithValue(ctx, ctxKey, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UUIDFromContext retrieves parsed uuid key from context.
// It implies that UUIDMiddleware was used
func UUIDFromContext(ctx context.Context, key string) uuid.UUID {
	ctxKey := uuidCtxKey{key: key}
	id, ok := ctx.Value(ctxKey).(uuid.UUID)
	if !ok {
		panic(fmt.Sprintf("uuid not found in context: %s. uuid middleware is not applied", key))
	}

	return id
}
