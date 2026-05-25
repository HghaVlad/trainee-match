package eventhandler

import "errors"

type Event struct {
	Topic   string
	Headers map[string][]byte
	Key     []byte
	Payload []byte
}

type ResultStatus string

const (
	ResultStatusSuccess ResultStatus = "success"
	ResultStatusRetry   ResultStatus = "error"
	ResultStatusDLQ     ResultStatus = "dlq"
)

var ErrUnknownEventType = errors.New("unknown event type")
