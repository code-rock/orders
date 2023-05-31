package order

import (
	"context"
	"log"
	"order-list/internal/database/postgresql"

	"github.com/jackc/pgconn"
)

func (r *SRepository) Create(ctx context.Context, order *SOrderTable) error {
	q := `
		INSERT INTO oreder_list
			(order_uid, bin)
		VALUES
			($1, $2)
		RETURNING order_uid		
	`
	if err := r.client.QueryRow(ctx, q, order.ID, order.Bin).Scan(&order.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			log.Printf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			return nil
		}
		return err
	}

	return nil
}

func (r *SRepository) FindAll(ctx context.Context) (u []SOrderTable, err error) {
	q := `SELECT order_uid, bin FROM oreder_list;`
	log.Printf("SQL Query: %s", q)

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	orders := make([]SOrderTable, 0)

	for rows.Next() {
		var curr SOrderTable

		err = rows.Scan(&curr.ID, &curr.Bin)
		if err != nil {
			return nil, err
		}

		orders = append(orders, curr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func NewRepository(client postgresql.SClient, logger interface{}) IRepository {
	return &SRepository{
		client: client,
	}
}
