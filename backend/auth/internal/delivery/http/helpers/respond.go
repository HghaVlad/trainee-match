package helpers

import (
	"encoding/json"
	"log/slog"

	"net/http"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/delivery/http/dto"
)

func RespondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("json encode error", "error", err)
	}
}

func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, dto.ErrorResponse{
		Error: message,
	})
}
