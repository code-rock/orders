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

func NewRepository(client postgresql.SClient, logger interface{}) order.Repository {
	return &Repository{
		client: client,
		logger: logger,
	}
}
