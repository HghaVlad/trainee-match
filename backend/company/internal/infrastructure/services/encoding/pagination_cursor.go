package encoding

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
)

type CursorWrapper[OrderT comparable] struct {
	Order OrderT
	Data  []byte
}

func DecodeCursor[T any, OrderT comparable](raw string, expectedOrder OrderT) (*T, error) {
	if raw == "" {
		return nil, nil
	}

	b, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, errors.Join(domain_errors.ErrInvalidCursor, err)
	}

	var wrapper CursorWrapper[OrderT]
	if err := json.Unmarshal(b, &wrapper); err != nil {
		return nil, errors.Join(domain_errors.ErrInvalidCursor, err)
	}

	if wrapper.Order != expectedOrder {
		return nil, domain_errors.ErrCursorOrderMismatch
	}

	var result T
	if err := json.Unmarshal(wrapper.Data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func EncodeCursor[T any, OrderT comparable](order OrderT, data *T) (*string, error) {
	if data == nil {
		return nil, nil
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	wrapper := CursorWrapper[OrderT]{
		Order: order,
		Data:  payload,
	}

	bytes, err := json.Marshal(wrapper)
	if err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString(bytes)
	return &encoded, nil
}
