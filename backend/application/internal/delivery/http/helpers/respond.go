package helpers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/HghaVlad/trainee-match/backend/application/internal/delivery/http/dto"
)

func RespondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Info("json encode error: %v", err)
	}
}

func RespondWithError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError

	switch {
	case errors.Is(err, dto.ErrBadRequest):
		code = http.StatusBadRequest
		break
	}
	RespondJSON(w, code, dto.JSONResponse{Message: err.Error()})
}
