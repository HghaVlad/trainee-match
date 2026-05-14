package eventhandler

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
