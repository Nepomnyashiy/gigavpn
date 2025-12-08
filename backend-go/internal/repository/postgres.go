package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Config для подключения к базе данных.
type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

// NewPostgresDB создает новый пул соединений с PostgreSQL.
func NewPostgresDB(ctx context.Context, cfg DBConfig) (*pgxpool.Pool, error) {
	// postgresql://user:password@host:port/dbname
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать пул соединений: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	return pool, nil
}
