package helpers

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/M0s1ck/g-store/src/pkg/http/responds"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list"
)

func ParseLimit(r *http.Request, key string, defaultLimit int) int {
	str := r.URL.Query().Get(key)
	limit, err := strconv.Atoi(str)

	if err != nil || limit <= 0 {
		limit = defaultLimit
	}

	return limit
}

func parseRangeInt(q url.Values, minKey, maxKey string) *list.RangeInt {
	var r list.RangeInt
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

func ParseUuidFromPathOr400(r *http.Request, w http.ResponseWriter, key string) (uuid.UUID, bool) {
	str := chi.URLParam(r, key)
	if str == "" {
		responds.RespondError(w, http.StatusBadRequest, errors.New("uuid parameter is required"))
		return uuid.Nil, false
	}

	val, err := uuid.Parse(str)
	if err != nil {
		responds.RespondError(w, http.StatusBadRequest, err)
		return uuid.Nil, false
	}

	return val, true
}
