package common

import "errors"

var ErrInvalidCursor = errors.New("invalid cursor")
var ErrCursorOrderMismatch = errors.New("cursor order mismatch")
var ErrUnsupportedListOrder = errors.New("unsupported list order")

var ErrLimitTooLarge = errors.New("limit too large")
