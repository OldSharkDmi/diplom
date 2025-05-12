package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("PGUSER"), os.Getenv("PGPASS"),
		os.Getenv("PGHOST"), os.Getenv("PGPORT"), os.Getenv("PGDB"),
	)
	return pgxpool.New(ctx, dsn)
}
