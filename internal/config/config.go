package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type ServerConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
	SSLMode  string `json:"ssl_modeMode"`
}

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	JWTSecret string `json:"jwt_secret"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	cfg := Config{}

	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return &cfg, nil
}

func ValidateConfig(cfg *Config) error {
	var errors []string

	if cfg.Server.Host == "" {
		errors = append(errors, "server_host обязателен")
	}

	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		errors = append(errors, "недопустимое значение для server_port")
	}

	if cfg.Database.Host == "" {
		errors = append(errors, "database_host обязателен")
	}

	if cfg.Database.Port < 1 || cfg.Database.Port > 65535 {
		errors = append(errors, "недопустимое значение для database_port")
	}

	if cfg.Database.User == "" {
		errors = append(errors, "database_user обязателен")
	}

	if len(errors) > 0 {
		return fmt.Errorf("ошибки валидации: %s", strings.Join(errors, "; "))
	}
	return nil
}

func ApplyDefaults(cfg *Config) {
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 30
	}

	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 30
	}

	if cfg.Database.Port == 0 {
		cfg.Database.Port = 5432
	}

	if cfg.Database.SSLMode == "" {
		cfg.Database.SSLMode = "disabled"
	}
}

func NewDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host:         "localhost",
			Port:         8080,
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		Database: DatabaseConfig{
			Host:    "localhost",
			Port:    5432,
			SSLMode: "disable",
		},
	}
}
