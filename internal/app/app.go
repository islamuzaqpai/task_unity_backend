package app

import (
	"enactus/internal/auth"
	"enactus/internal/config"
	"enactus/internal/database"
	"enactus/internal/handler"
	"enactus/internal/repository"
	"enactus/internal/routes"
	"enactus/internal/service"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
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

	userRepo := repository.UserRepository{Pool: pool}
	userS := service.UserService{
		UserRepo:  &userRepo,
		JwtSecret: &auth.JWTSecret{Secret: []byte(cfg.JWTSecret)},
	}
	userH := handler.UserHandler{UserS: &userS}

	router := gin.Default()
	routes.UserRoutes(&userH, router)

	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("fail: %v", err)
	}
}
