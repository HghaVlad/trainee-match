package list_companies

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	domain_errors "github.com/HghaVlad/trainee-match/backend/company/internal/domain/errors"
)

type CursorWrapper struct {
	Order Order
	Data  []byte
}

// Maybe move to somewhere common for reuse

func decodeCursor[T any](raw string, expectedOrder Order) (*T, error) {
	if raw == "" {
		return nil, nil
	}

	b, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, domain_errors.ErrInvalidCursor)
	}

	var wrapper CursorWrapper
	if err := json.Unmarshal(b, &wrapper); err != nil {
		return nil, fmt.Errorf("%v: %w", domain_errors.ErrInvalidCursor, err)
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

func encodeCursor[T any](order Order, data *T) (*string, error) {
	if data == nil {
		return nil, nil
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	wrapper := CursorWrapper{
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
