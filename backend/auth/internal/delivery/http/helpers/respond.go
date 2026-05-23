package helpers

import (
	"encoding/json"
	"errors"
	"log/slog"

	"net/http"

	"github.com/HghaVlad/trainee-match/backend/auth/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/auth/internal/domain"
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

func RespondSmartError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, domain.ErrInvalidName) || errors.Is(err, domain.ErrInvalidEmail) || errors.Is(err, domain.ErrInvalidRole):
		status = http.StatusBadRequest
	case errors.Is(err, domain.ErrEmailAlreadyExists) || errors.Is(err, domain.ErrUserNameAlreadyExists):
		status = http.StatusConflict
	case errors.Is(err, domain.ErrIncorrectPassword) || errors.Is(err, domain.ErrUnauthorized):
		status = http.StatusUnauthorized
	}
	RespondError(w, status, err.Error())
}
