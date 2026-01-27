package update_company

import "context"

type CompanyRepo interface {
	Update(ctx context.Context, req *Request) error
}
