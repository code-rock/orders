package order

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, order *SOrderTable) error
	FindAll(ctx context.Context) (u []SOrderTable, err error)
	// clearTable(ctx context.Context)
	// FindOne(ctx context.Context, id string) (SOrderTable, error)
	// Update(ctx context.Context, order SOrderTable) error
	// Delete(ctx context.Context, id string) error
}
