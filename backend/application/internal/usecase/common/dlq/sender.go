package dlq

import "context"

type Producer interface {
	ProduceDLQ(ctx context.Context, message Message, key, value []byte) error
}

type Encoder interface {
	DLQToBytes(message Message) ([]byte, error)
}

type Sender struct {
	producer Producer
	encoder  Encoder
}

func NewSender(producer Producer, encoder Encoder) *Sender {
	return &Sender{producer, encoder}
}

func (s *Sender) SendDLQ(ctx context.Context, message Message, key []byte) error {
	value, err := s.encoder.DLQToBytes(message)
	if err != nil {
		return err
	}
	return s.producer.ProduceDLQ(ctx, message, key, value)
}
