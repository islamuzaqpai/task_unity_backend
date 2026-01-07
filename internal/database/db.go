package database

import (
	"context"
	"enactus/internal/config"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(dbCfg config.DB) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s",
		dbCfg.Host,
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Port,
		dbCfg.DBName,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to configure pool: %w", err)
	}

	poolCfg.MinConns = 2
	poolCfg.HealthCheckPeriod = time.Minute * 5
	poolCfg.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %v", err)
	}

	return pool, nil
}
