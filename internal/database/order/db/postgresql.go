package dbOrder

import (
	"basket/internal/database/order"
	"basket/internal/database/postgresql"
	"context"
	"fmt"

	"github.com/jackc/pgconn"
)

type Repository struct {
	client postgresql.SClient
	logger interface{} //*logging.Logger
}

func (r *Repository) Create(ctx context.Context, order *order.SOrderTable) error {
	q := `
		INSERT INTO oreder_list
			(order_uid, bin)
		VALUES
			($1, $2)
		RETURNING order_uid		
	`
	if err := r.client.QueryRow(ctx, q, order.ID, order.Bin).Scan(&order.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			fmt.Println(newErr)
			return nil
		}
		return err
	}

	return nil
}

func (r *Repository) FindAll(ctx context.Context) (u []order.SOrderTable, err error) {
	q := `SELECT order_uid, bin FROM oreder_list;`
	// r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	orders := make([]order.SOrderTable, 0)

	for rows.Next() {
		var curr order.SOrderTable

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

func NewRepository(client postgresql.SClient, logger interface{}) order.Repository {
	return &Repository{
		client: client,
		logger: logger,
	}
}
