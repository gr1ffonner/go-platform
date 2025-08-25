package mysql

import (
	"context"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type MySQLClient struct {
	db *sqlx.DB
}

func NewMySQL(ctx context.Context, dsn string) (*MySQLClient, error) {
	// Parse DSN to validate it
	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	db, err := sqlx.ConnectContext(ctx, "mysql", config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping mysql: %w", err)
	}

	return &MySQLClient{db: db}, nil
}

func (m *MySQLClient) Close() {
	if m.db != nil {
		m.db.Close()
	}
}

func (m *MySQLClient) DB() *sqlx.DB {
	return m.db
}
