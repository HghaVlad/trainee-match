package listbycomp

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"slices"
	"sort"
	"time"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	vaclist "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list"
)

type ResponseCacheRepo interface {
	Get(ctx context.Context, key string) *Response
	Put(ctx context.Context, key string, response *Response, exp time.Duration)
}

type cacheReq struct {
	CompID       string                `json:"compId"`
	Order        Order                 `json:"order"`
	Limit        int                   `json:"limit"`
	Cursor       string                `json:"cursor"`
	Requirements *vaclist.Requirements `json:"requirements,omitempty"`
	Status       *vacancy.Status       `json:"status,omitempty"`
}

func requestToCacheKey(r *Request) string {
	payload := cacheReq{
		CompID:       r.CompID.String(),
		Order:        r.Order,
		Limit:        r.Limit,
		Cursor:       r.EncodedCursor,
		Requirements: normalizeRequirements(r.Requirements),
		Status:       r.Status,
	}

	//nolint:musttag // marshal to inner cache, unmarshal by the same rules
	b, _ := json.Marshal(payload)

	sum := sha256.Sum256(b)

	return hex.EncodeToString(sum[:])
}

func normalizeRequirements(req *vaclist.Requirements) *vaclist.Requirements {
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
