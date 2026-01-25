package helpers

import (
	"net/http"
	"strconv"
)

func ParseLimit(r *http.Request, key string, defaultLimit int) int {
	str := r.URL.Query().Get(key)
	limit, err := strconv.Atoi(str)

	if err != nil || limit <= 0 {
		limit = defaultLimit
	}

	return limit
}
