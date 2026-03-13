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

var ErrHrRoleRequired = errors.New("hr role is required")
var ErrInsufficientRole = errors.New("insufficient role")

var ErrCompanyMemberRequired = errors.New("being this company's member is required")
var ErrInsufficientRoleInCompany = errors.New("insufficient company member role")

var ErrCompanyMemberNotFound = errors.New("company member not found")
var ErrCompanyMemberAlreadyExists = errors.New("company member already exists")
var ErrInvalidUserID = errors.New("invalid user id")
var ErrInvalidCompanyMemberRole = errors.New("invalid company member role")

var ErrEmptyCompaniesFilter = errors.New("empty companies filter")
var ErrEmptyCityFilter = errors.New("empty city filter")

var ErrLimitTooLarge = errors.New("limit too large")

var ErrInvalidSalaryOrderForUnpaid = errors.New("unpaid vacancies don't support this order")
