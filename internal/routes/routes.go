package routes

import (
	"enactus/internal/handler"
	"enactus/internal/httpx"
	"net/http"
)

func UserRoutes(userH *handler.UserHandler, mux *http.ServeMux) {
	mux.HandleFunc("POST /users", httpx.WrapHandler(userH.Register))
	mux.HandleFunc("GET /users", httpx.WrapHandler(userH.GetAllUsers))
}
