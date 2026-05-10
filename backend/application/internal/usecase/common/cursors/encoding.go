package cursors

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

type summaryCursorWrapper struct {
	Order SummaryOrder    `json:"order"`
	Data  json.RawMessage `json:"data"`
}

type hrSummaryCursorWrapper struct {
	Order HrSummaryOrder  `json:"order"`
	Data  json.RawMessage `json:"data"`
}

func DecodeSummaryCursor(raw string, expectedOrder SummaryOrder) (*SummaryCursor, error) {
	if raw == "" {
		//nolint:nilnil // nil cursor is possible
		return nil, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, errors.Join(ErrInvalidCursor, err)
	}

	var wrapper summaryCursorWrapper
	if err := json.Unmarshal(decoded, &wrapper); err != nil {
		return nil, errors.Join(ErrInvalidCursor, err)
	}

	if wrapper.Order != expectedOrder {
		return nil, ErrCursorOrderMismatch
	}

	var cursor SummaryCursor
	if err := json.Unmarshal(wrapper.Data, &cursor); err != nil {
		return nil, errors.Join(ErrInvalidCursor, err)
	}

	return &cursor, nil
}

func EncodeCursor(order SummaryOrder, cursor *SummaryCursor) (*string, error) {
	if cursor == nil {
		//nolint:nilnil // nil cursor is possible
		return nil, nil
	}

	payload, err := json.Marshal(cursor)
	if err != nil {
		return nil, err
	}

	wrapper, err := json.Marshal(summaryCursorWrapper{
		Order: order,
		Data:  payload,
	})
	if err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString(wrapper)
	return &encoded, nil
}

func DecodeHrSummaryCursor(raw string, expectedOrder HrSummaryOrder) (*HrSummaryCursor, error) {
	if raw == "" {
		//nolint:nilnil // nil cursor is possible
		return nil, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, errors.Join(ErrInvalidCursor, err)
	}

	var wrapper hrSummaryCursorWrapper
	if err := json.Unmarshal(decoded, &wrapper); err != nil {
		return nil, errors.Join(ErrInvalidCursor, err)
	}

	if wrapper.Order != expectedOrder {
		return nil, ErrCursorOrderMismatch
	}

	var cursor HrSummaryCursor
	if err := json.Unmarshal(wrapper.Data, &cursor); err != nil {
		return nil, errors.Join(ErrInvalidCursor, err)
	}

	return &cursor, nil
}

func EncodeHrSummaryCursor(order HrSummaryOrder, cursor *HrSummaryCursor) (*string, error) {
	if cursor == nil {
		//nolint:nilnil // nil cursor is possible
		return nil, nil
	}

	payload, err := json.Marshal(cursor)
	if err != nil {
		return nil, err
	}

	wrapper, err := json.Marshal(hrSummaryCursorWrapper{
		Order: order,
		Data:  payload,
	})
	if err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString(wrapper)
	return &encoded, nil
}
