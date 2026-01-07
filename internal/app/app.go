package app

import (
	"enactus/internal/config"
	"enactus/internal/database"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func Run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("fail: %v", err)
	}

	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("failed to load data from .env: %v", err)
	}

	pool, err := database.Connect(cfg.DB)
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}

	defer pool.Close()

	fmt.Println("Success", pool)
}
