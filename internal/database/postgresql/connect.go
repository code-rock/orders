package postgresql

import (
	envconfig "basket/internal/config"
	"basket/internal/utils"
	"fmt"
	"time"

	// "database/sql"
	// "fmt"
	// "log"
	// "context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"

	"context"
)

//	type Client interface {
//		Ex
//	}
//
// orders=# CREATE TABLE oreder_list (
//
//	order_uid varchar(20) PRIMARY KEY,
//	bin bytea NOT NULL
//	);
type SClient interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, maxAttempts int, params envconfig.SDBConfig) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", params.User, params.Password, params.Host, params.Port, params.DBName)

	utils.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			fmt.Println("sfsefe")
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	return pool, nil
}

// func Connect(params envconfig.SDBConfig) {
// 	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
// 		params.Host, params.Port, params.User, params.Password, params.DBName)
// 	// // "user=tania dbname=test sslmode=require host=localhost"
// 	db, err := sql.Open("postgres", connStr)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("db")
// 	fmt.Println(db)
// 	fmt.Println("db")
// 	// age := 21
// 	// rows, err := db.Query("SELECT name FROM user WHERE age = $1", age)

// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }
// 	// fmt.Println(rows)
// }
