package routes

import (
	"enactus/internal/auth"
	"enactus/internal/handler"
	"enactus/internal/httpx"
	"enactus/internal/middleware"
	"enactus/internal/models"
	"net/http"
)

func UserRoutes(userH *handler.UserHandler, mux *http.ServeMux, jwtSecret *auth.JWTSecret) {
	mux.HandleFunc("PATCH /users/update/profile/{id}", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret, userH.UpdateUserProfile)))
	mux.HandleFunc("PATCH /users/update/password/{id}", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret, userH.UpdateUserPassword)))
	mux.HandleFunc("DELETE /users/delete/{id}", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(userH.DeleteUser, models.Role{Name: "admin"}))))

	mux.HandleFunc("POST /users/register", httpx.WrapHandler(userH.Register))
	mux.HandleFunc("GET /users", httpx.WrapHandler(userH.GetAllUsers))
	mux.HandleFunc("GET /users/{id}", httpx.WrapHandler(userH.GetUserById))
	mux.HandleFunc("POST /users/login", httpx.WrapHandler(userH.Login))

}

func TaskRoutes(taskH *handler.TaskHandler, mux *http.ServeMux, jwtSecret *auth.JWTSecret) {
	mux.HandleFunc("POST /tasks/create", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(taskH.AddTask, models.Role{Name: "admin"}, models.Role{Name: "manager"}))))
	mux.HandleFunc("GET /tasks/assignee", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret, taskH.GetAllTasksByAssigneeId)))
	mux.HandleFunc("PATCH /tasks/update/{id}", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(taskH.UpdateTask, models.Role{Name: "admin"}, models.Role{Name: "manager"}))))
	mux.HandleFunc("DELETE /tasks/delete/{id}", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(taskH.DeleteTask, models.Role{Name: "admin"}, models.Role{Name: "manager"}))))

	mux.HandleFunc("GET /tasks/{id}", httpx.WrapHandler(taskH.GetTaskById))
	mux.HandleFunc("GET /tasks/", httpx.WrapHandler(taskH.GetAllTasks))
}

func DepartmentRoutes(departmentH *handler.DepartmentHandler, mux *http.ServeMux, jwtSecret *auth.JWTSecret) {
	mux.HandleFunc("POST /departments/create", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(departmentH.AddDepartment, models.Role{Name: "admin"}))))
	mux.HandleFunc("GET /departments/", httpx.WrapHandler(departmentH.GetAllDepartments))
}

func AttendanceRoutes(attendanceH *handler.AttendanceHandler, mux *http.ServeMux, jwtSecret *auth.JWTSecret) {
	mux.HandleFunc("POST /attendances/create", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(attendanceH.AddAttendance, models.Role{Name: "admin"}, models.Role{Name: "manager"}))))

	mux.HandleFunc("GET /attendances", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(attendanceH.GetAllAttendances, models.Role{Name: "admin"}, models.Role{Name: "manager"}))))

	mux.HandleFunc("PATCH /attendances/update/{id}", httpx.WrapHandler(middleware.AuthMiddleware(jwtSecret,
		middleware.RoleMiddleware(attendanceH.UpdateAttendance, models.Role{Name: "admin"}, models.Role{Name: "manager"}))))
}
