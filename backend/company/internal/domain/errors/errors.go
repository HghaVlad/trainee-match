package domain_errors

import "errors"

var ErrCompanyNotFound = errors.New("company not found")

var ErrCompanyAlreadyExists = errors.New("company with this name already exists")

var ErrInvalidCursor = errors.New("invalid cursor")

var ErrCursorOrderMismatch = errors.New("cursor order mismatch")

var ErrUnsupportedListOrder = errors.New("unsupported list order")

var ErrVacancyNotFound = errors.New("vacancy not found")
