package listbycomp

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type Request struct {
	CompID uuid.UUID
	Order  Order
	Cursor string
	Limit  int
}

func (r *Request) toCacheKey() string {
	return strings.Join([]string{
		r.CompID.String(), string(r.Order), r.Cursor, strconv.Itoa(r.Limit)}, "-")
}
