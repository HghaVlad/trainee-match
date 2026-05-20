package dlq

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	EventID           uuid.UUID `avro:"event_id"`
	Payload           []byte    `avro:"payload"`
	OriginalTopic     string    `avro:"original_topic"`
	OriginalEventType string    `avro:"original_event_type"`
	LastError         string    `avro:"last_error"`
	FailedAt          time.Time `avro:"failed_at"`
}
