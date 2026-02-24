package list_vac_by_comp

import (
	"time"

	"github.com/google/uuid"
)

type Order string

const (
	OrderPublishedAtDesc Order = "published_at_desc"
)

type PublishedAtCursor struct {
	PublishedAt time.Time
	Id          uuid.UUID
}
