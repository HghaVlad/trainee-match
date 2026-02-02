package domain_errors

import "errors"

var ErrCompanyNotFound = errors.New("company not found")

var ErrCompanyAlreadyExists = errors.New("company with this name already exists")

var ErrCompanyInvalidNameLen = errors.New("invalid company name length")
var ErrCompanyInvalidDescriptionLen = errors.New("invalid company description length")

var ErrInvalidCursor = errors.New("invalid cursor")
var ErrCursorOrderMismatch = errors.New("cursor order mismatch")
var ErrUnsupportedListOrder = errors.New("unsupported list order")

var ErrVacancyNotFound = errors.New("vacancy not found")

var ErrInvalidWorkFormat = errors.New("invalid work format")
var ErrInvalidEmploymentType = errors.New("invalid employment type")

var ErrInvalidDurationRange = errors.New("invalid duration range")
var ErrInvalidHoursRange = errors.New("invalid hours per week range")
var ErrInvalidSalaryRange = errors.New("invalid salary range")

var ErrSalaryProvidedForUnpaid = errors.New("salary provided for unpaid vacancy")
var ErrSalaryMissingForPaid = errors.New("salary missing for paid vacancy")

var ErrInvalidTitleLength = errors.New("invalid title length")
var ErrInvalidDescriptionLength = errors.New("invalid description length")

var ErrSalaryTooLarge = errors.New("salary too large")
var ErrNegativeSalary = errors.New("salary negative")
