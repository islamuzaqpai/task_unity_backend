package routes

import (
	"enactus/internal/auth"
	"enactus/internal/handler"
	"enactus/internal/httpx"
	"enactus/internal/middleware"
	"net/http"
)

func UserRoutes(userH *handler.UserHandler, mux *http.ServeMux, jwtSecret *auth.JWTSecret) {
	mux.HandleFunc("POST /users/register", httpx.WrapHandler(userH.Register))
	mux.HandleFunc("GET /users", httpx.WrapHandler(userH.GetAllUsers))
	mux.HandleFunc("GET /users/{id}", httpx.WrapHandler(userH.GetUserById))
	mux.HandleFunc("POST /users/login", httpx.WrapHandler(userH.Login))
	mux.HandleFunc("DELETE /users/delete/{id}", httpx.WrapHandler(userH.DeleteUser))

	//with middleware
	mux.HandleFunc("PATCH /users/update/profile/{id}", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret, userH.UpdateUserProfile)))
	mux.HandleFunc("PATCH /users/update/password/{id}", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret, userH.UpdateUserPassword)))

}

func TaskRoutes(taskH *handler.TaskHandler, mux *http.ServeMux, jwtSecret *auth.JWTSecret) {
	mux.HandleFunc("POST /tasks/create", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret, taskH.AddTask)))
	mux.HandleFunc("GET /tasks/assignee", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret, taskH.GetAllTasksByAssigneeId)))
	mux.HandleFunc("PATCH /tasks/update/{id}", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret, taskH.UpdateTask)))

	mux.HandleFunc("GET /tasks/{id}", httpx.WrapHandler(taskH.GetTaskById))
	mux.HandleFunc("GET /tasks/", httpx.WrapHandler(taskH.GetAllTasks))
}
