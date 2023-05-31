package postgresql

import (
	"fmt"
	envconfig "order-list/internal/config"
	"order-list/internal/utils"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"

	"context"
)

type SClient interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func NewClient(ctx context.Context, maxAttempts int, params envconfig.SDBConfig) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", params.User, params.Password, params.Host, params.Port, params.DBName)

	utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	return pool, nil
}
