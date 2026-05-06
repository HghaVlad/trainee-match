package member

import "errors"

var ErrCompanyMemberNotFound = errors.New("company member not found")
var ErrCompanyMemberAlreadyExists = errors.New("company member already exists")
var ErrInvalidUserID = errors.New("invalid user id")
var ErrInvalidCompanyMemberRole = errors.New("invalid company member role")

var ErrCompanyMemberRequired = errors.New("being this company's member is required")

var ErrInsufficientRoleInCompany = errors.New("insufficient company member role")

var ErrCantRemoveYourself = errors.New("can't remove yourself")
