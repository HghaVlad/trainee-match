package helpers

import (
	"context"
	"encoding/json"
	"net/http"

	utilslog "github.com/HghaVlad/trainee-match/backend/company/internal/infrastructure/utils/logger"
	"github.com/HghaVlad/trainee-match/backend/company/internal/transport/http/dto"
)

func RespondJSON(ctx context.Context, w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		logger := utilslog.FromContext(ctx)
		logger.ErrorContext(ctx, "json encode error", "err", err)
	}
}

func RespondError(ctx context.Context, w http.ResponseWriter, status int, err error) {
	RespondJSON(ctx, w, status, dto.ErrorResponse{
		Error: err.Error(),
	})
}

func RespondErrorMsg(ctx context.Context, w http.ResponseWriter, status int, msg string) {
	RespondJSON(ctx, w, status, dto.ErrorResponse{
		Error: msg,
	})
}
