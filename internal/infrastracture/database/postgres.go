package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresConfig struct {
	PostgresUser, PostgresPass, PostgresName, PostgresHost, PostgresPort string
}

func GetPostgresDatabase(ctx context.Context, config PostgresConfig) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.PostgresUser, config.PostgresPass, config.PostgresHost, config.PostgresPort, config.PostgresName,
	)
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("error while connecting to database: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	return pool, nil
}
