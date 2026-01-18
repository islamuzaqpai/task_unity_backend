package app

import (
	"context"
	"enactus/internal/config"
	"enactus/internal/database"
	"enactus/internal/repository"
	"enactus/internal/service"
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

	commentRepo := repository.CommentRepository{Pool: pool}

	commentS := service.CommentService{CommentRepo: &commentRepo}

	err = commentS.DeleteComment(context.Background(), 5)
	if err != nil {
		log.Fatalf("fail: %v", err)
	}

	fmt.Println("yes")
}
