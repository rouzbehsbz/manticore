package infra

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rouzbehsbz/manticore/server/internal/infra/db/sources"
)

type Db struct {
	Q *sources.Queries
}

func NewPostgresService(host string, port uint16, username, password, databaseName string, maxConnections int) (*Db, error) {
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", username, password, host, port, databaseName)

	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	config.MaxConns = int32(maxConnections)

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %v", err)
	}

	q := sources.New(pool)

	return &Db{
		Q: q,
	}, nil
}
