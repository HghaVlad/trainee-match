package helpers

import (
	"context"
	"encoding/json"
	"net/http"

	utilslog "github.com/HghaVlad/trainee-match/backend/application/internal/infrastructure/utils/logger"
)

type ErrorResponse struct {
	Error   string         `json:"error"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

func WriteError(
	ctx context.Context,
	w http.ResponseWriter,
	status int,
	errCode string,
	message string,
	details map[string]any,
) {
	resp := ErrorResponse{
		Error:   errCode,
		Message: message,
		Details: details,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger := utilslog.FromContext(ctx)
		logger.ErrorContext(ctx, "json encode error", "err", err)
	}
}
