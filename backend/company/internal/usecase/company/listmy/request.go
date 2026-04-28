package listmy

import "github.com/HghaVlad/trainee-match/backend/company/internal/usecase/company/list"

type Request struct {
	Order         list.Order
	Limit         int
	EncodedCursor string
}
