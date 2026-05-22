package helpers

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/listsearch"
)

func ParseLimit(r *http.Request, key string, defaultLimit int) int {
	str := r.URL.Query().Get(key)
	limit, err := strconv.Atoi(str)

	if err != nil || limit <= 0 {
		limit = defaultLimit
	}

	return limit
}

func parseRangeInt(q url.Values, minKey, maxKey string) *listsearch.RangeInt {
	var r listsearch.RangeInt
	var hasValue bool

	if minStr := q.Get(minKey); minStr != "" {
		if mn, err := strconv.Atoi(minStr); err == nil {
			r.Min = &mn
			hasValue = true
		}
	}

	if maxStr := q.Get(maxKey); maxStr != "" {
		if mx, err := strconv.Atoi(maxStr); err == nil {
			r.Max = &mx
			hasValue = true
		}
	}

	if !hasValue {
		return nil
	}

	return &r
}
