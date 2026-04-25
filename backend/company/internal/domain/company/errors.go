package company

import "errors"

var ErrCompanyNotFound = errors.New("company not found")

var ErrCompanyAlreadyExists = errors.New("company with this name already exists")

var ErrCompanyInvalidNameLen = errors.New("invalid company name length")
var ErrCompanyInvalidDescriptionLen = errors.New("invalid company description length")
