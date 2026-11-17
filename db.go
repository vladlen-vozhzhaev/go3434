package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func initDB() error {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		connString = "postgres://postgres:qwerty123@localhost:5432/blogdb"
	}

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("could not ping database: %w", err)
	}

	db = pool
	return nil
}
