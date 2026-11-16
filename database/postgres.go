package database

import (
	"context"
	"fmt"
	"time"

	"github.com/AanFLenrir/PetrivaBackend/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}
	// Try ping
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("cannot connect to database: %w", err)
	}
	return pool, nil
}
