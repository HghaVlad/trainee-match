package eventhandler

type Event struct {
	Topic   string
	Headers map[string][]byte
	Key     []byte
	Payload []byte
}
