package app

import (
	"enactus/internal/config"
	"enactus/internal/database"
	"enactus/internal/handler"
	"enactus/internal/repository"
	"enactus/internal/routes"
	"enactus/internal/service"
	"fmt"
	"log"
	"net/http"

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

	userR := repository.NewUserRepository(pool)
	userS := service.NewUserService(userR)
	userH := handler.NewUserHandler(userS)

	mux := http.NewServeMux()
	routes.UserRoutes(userH, mux)

	addr := ":8080"
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
