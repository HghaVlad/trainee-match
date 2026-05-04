package eventhandler

import "errors"

type Event struct {
	Payload []byte
	Key     []byte
	Topic   string
	Headers map[string][]byte
}

type ResultStatus int

const (
	ResultSuccess ResultStatus = iota
	ResultRetry
	ResultDLQ
)

var ErrUnknownEventType = errors.New("unknown event type")

const UserCreatedEventType = "UserCreated"
