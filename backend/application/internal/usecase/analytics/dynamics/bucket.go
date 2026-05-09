package dynamics

import (
	"time"
)

type Bucket struct {
	Start          time.Time
	SubmittedCount int
	InterviewCount int
	OfferCount     int
	RejectedCount  int
	SeenCount      int
	WithdrawnCount int
}
