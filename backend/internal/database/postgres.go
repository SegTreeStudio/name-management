package database

import (
    "context"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {

    config, err := pgxpool.ParseConfig(databaseURL)
    if err != nil {
        return nil, fmt.Errorf("parse database config: %w", err)
    }

    config.MaxConns = 10
    config.MinConns = 2

    ctx, cancel := context.WithTimeout(
        context.Background(),
        10*time.Second,
    )
    defer cancel()

    pool, err := pgxpool.NewWithConfig(ctx, config)
    if err != nil {
        return nil, fmt.Errorf("create connection pool: %w", err)
    }

    // Actually verify the database connection.
    if err := pool.Ping(ctx); err != nil {
        pool.Close()
        return nil, fmt.Errorf("database ping failed: %w", err)
    }

    return pool, nil
}