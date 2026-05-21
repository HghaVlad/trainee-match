package helpers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/go-chi/chi/v5"
)

func ParseUUIDParam(r *http.Request, param string) (uuid.UUID, error) {
	raw := chi.URLParam(r, param)
	if raw == "" {
		return uuid.Nil, errors.New("id is required")
	}

	parsed, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, errors.New("invalid id format")
	}

	return parsed, nil
}

func ParsePageSize(r *http.Request) (int, int, error) {
	query := r.URL.Query()

	page := 1
	size := 20

	if raw := query.Get("page"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value <= 0 {
			return 0, 0, errors.New("invalid page")
		}
		page = value
	}

	if raw := query.Get("size"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value <= 0 {
			return 0, 0, errors.New("invalid size")
		}
		size = value
	}

	return page, size, nil
}
