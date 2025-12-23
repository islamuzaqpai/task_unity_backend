package database

import (
	"context"
	"enactus/internal/config"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"time"
)

func Connect() (*pgxpool.Pool, error) {
	err := godotenv.Load(".env") //убрать (нужно чтобы эта функция вызывалась 1 раз)
	if err != nil {
		log.Fatalf("error loading env file: %v", err)
	}

	cfg, err := config.LoadFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to load data from .env:%w", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s",
		cfg.DB.Host,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Port,
		cfg.DB.DBName,
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
