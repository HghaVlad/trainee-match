package vacancy

import "errors"

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
var ErrInvalidStatus = errors.New("invalid vacancy status")

var ErrSalaryTooLarge = errors.New("salary too large")
var ErrNegativeSalary = errors.New("salary negative")

var ErrEmptyCompaniesFilter = errors.New("empty companies filter")
var ErrEmptyCityFilter = errors.New("empty city filter")

var ErrInvalidSalaryOrderForUnpaid = errors.New("unpaid vacancies don't support this order")
