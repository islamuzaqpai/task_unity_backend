package app

import (
	"enactus/internal/auth"
	"enactus/internal/config"
	"enactus/internal/database"
	"enactus/internal/handler"
	"enactus/internal/middleware"
	"enactus/internal/repository"
	"enactus/internal/routes"
	"enactus/internal/service"
	"fmt"
	"log"
	"net/http"
)

func Run() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("failed to load data from .env: %v", err)
	}

	err = config.ValidateConfig(cfg)
	if err != nil {
		log.Fatalf("invalid token: %v", err)
	}

	pool, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}

	defer pool.Close()

	fmt.Println("Success", pool)

	jwtSecret := auth.JWTSecret{Secret: []byte(cfg.JWTSecret)}
	userR := repository.NewUserRepository(pool)
	userS := service.NewUserService(userR, &jwtSecret)
	userH := handler.NewUserHandler(userS)

	taskR := repository.NewTaskRepo(pool)
	taskS := service.NewTaskService(taskR, userR)
	taskH := handler.NewTaskHandler(taskS)

	departmentR := repository.NewDepartmentRepository(pool)
	departmentS := service.NewDepartmentService(departmentR)
	departmentH := handler.NewDepartmentHandler(departmentS)

	mux := http.NewServeMux()
	routes.UserRoutes(userH, mux, &jwtSecret)
	routes.TaskRoutes(taskH, mux, &jwtSecret)
	routes.DepartmentRoutes(departmentH, mux, &jwtSecret)

	muxWithCors := middleware.CORS(mux)
	addr := ":8080"
	server := &http.Server{
		Addr:    addr,
		Handler: muxWithCors,
	}

	log.Fatal(server.ListenAndServe())
}
