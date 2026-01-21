package domain_errors

import "errors"

var ErrCompanyNotFound = errors.New("company not found")

var ErrCompanyAlreadyExists = errors.New("company with this name already exists")
