package common

import "context"

type TxManager interface {
	Do(ctx context.Context, fn func(context.Context) error) error
}
