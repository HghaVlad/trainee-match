package list

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"slices"
	"sort"
	"time"
)

type ResponseCacheRepo interface {
	Get(ctx context.Context, key string) *Response
	Put(ctx context.Context, key string, response *Response, exp time.Duration)
}

type cacheReq struct {
	Order        Order         `json:"order"`
	Limit        int           `json:"limit"`
	Cursor       string        `json:"cursor"`
	Requirements *Requirements `json:"requirements,omitempty"`
}

func requestToCacheKey(r *Request) string {
	normalized := normalizeRequirements(r.Requirements)

	payload := cacheReq{
		Order:        r.Order,
		Limit:        r.Limit,
		Cursor:       r.EncodedCursor,
		Requirements: normalized,
	}

	b, _ := json.Marshal(payload)

	sum := sha256.Sum256(b)

	return hex.EncodeToString(sum[:])
}

func normalizeRequirements(req *Requirements) *Requirements {
	if req == nil {
		return nil
	}

	clone := *req

	if clone.WorkFormat != nil {
		s := slices.Clone(*clone.WorkFormat)
		slices.Sort(s)
		clone.WorkFormat = &s
	}

	if clone.City != nil {
		s := slices.Clone(*clone.City)
		slices.Sort(s)
		clone.City = &s
	}

	if clone.Companies != nil {
		s := slices.Clone(*clone.Companies)
		sort.Slice(s, func(i, j int) bool {
			return s[i].String() < s[j].String()
		})
		clone.Companies = &s
	}

	return &clone
}
