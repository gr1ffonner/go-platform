package clickhouse

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickHouseClient struct {
	conn driver.Conn
}

func NewClickHouse(ctx context.Context, dsn string) (*ClickHouseClient, error) {

	options, err := clickhouse.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ClickHouse DSN: %w", err)
	}

	conn, err := clickhouse.Open(options)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	return &ClickHouseClient{conn: conn}, nil
}

func (c *ClickHouseClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *ClickHouseClient) Conn() driver.Conn {
	return c.conn
}
