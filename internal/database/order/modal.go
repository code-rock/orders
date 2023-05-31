package order

import (
	"context"
	"order-list/internal/database/postgresql"
)

type SOrderTable struct {
	ID  string `json:"order_uid"`
	Bin []byte `json:"bin"`
}

type SRepository struct {
	client postgresql.SClient
}

type IRepository interface {
	Create(ctx context.Context, order *SOrderTable) error
	FindAll(ctx context.Context) (u []SOrderTable, err error)
}
