package config

import (
	"fmt"
	"os"
)

type Config struct {
	DB        DB
	JWTSecret string
}

type DB struct {
	Host     string
	User     string
	Password string
	Port     string
	DBName   string
}

func LoadFromEnv() (*Config, error) {
	cfg := &Config{
		DB: DB{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Port:     os.Getenv("DB_PORT"),
			DBName:   os.Getenv("DB_DATABASE"),
		},
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	if cfg.DB.Host == "" || cfg.DB.User == "" || cfg.DB.Password == "" || cfg.DB.Port == "" || cfg.DB.DBName == "" || cfg.JWTSecret == "" {
		return nil, fmt.Errorf("missing required env variable")
	}

	return cfg, nil
}
