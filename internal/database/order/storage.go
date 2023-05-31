package order

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, order *SOrderTable) error
	FindAll(ctx context.Context) (u []SOrderTable, err error)
}
