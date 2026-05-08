package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"

	"github.com/HghaVlad/trainee-match/backend/application/internal/infrastructure/utils/logger"
	"github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/helpers"
)

func LoggingMiddleware(next nethttp.StrictHTTPHandlerFunc, _ string) nethttp.StrictHTTPHandlerFunc {
	return func(
		ctx context.Context,
		w http.ResponseWriter,
		r *http.Request,
		request any,
	) (response any, err error) {
		resp, err := next(ctx, w, r, request)
		if err != nil {
			lgr := logger.FromContext(ctx)
			lgr.ErrorContext(
				ctx,
				"handler failed",
				slog.Any("err", err),
			)

			return resp, helpers.InternalServerError
		}

		return resp, err
	}
}
