package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitPostgres(dsn string) *pgxpool.Pool {
	if dsn == "" {
		_, err := fmt.Fprintln(os.Stderr, "No DATABASE_DSN provided")
		if err != nil {
			return &pgxpool.Pool{}
		}
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Failed to connect to PostgreSQL: %v\n", err)
		if err != nil {
			return &pgxpool.Pool{}
		}
		os.Exit(1)
	}

	if err := pool.Ping(ctx); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "PostgreSQL ping failed: %v\n", err)
		if err != nil {
			return &pgxpool.Pool{}
		}
		os.Exit(1)
	}

	Pool = pool
	return pool
}
