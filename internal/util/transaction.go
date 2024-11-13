package util

import (
	"context"
)

type Transaction interface{}

type TransactionFactory interface {
	WithTransaction(ctx context.Context, f func(Transaction) error) error
}
