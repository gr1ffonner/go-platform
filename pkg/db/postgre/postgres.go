package postgre

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClient struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, dsn string) (*PostgresClient, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return &PostgresClient{pool: pool}, nil
}

func (p *PostgresClient) Close() {
	if p.pool != nil {
		p.pool.Close()
	}
}

func (p *PostgresClient) Pool() *pgxpool.Pool {
	return p.pool
}
