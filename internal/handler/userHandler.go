package handler

import (
	"enactus/internal/httpx"
	"enactus/internal/models"
	"enactus/internal/service"
	"encoding/json"
	"log"
	"net/http"
)

type UserHandlerInterface interface {
	Register(w http.ResponseWriter, r *http.Request) error
	GetAllUsers(w http.ResponseWriter, r *http.Request) error
	Login(w http.ResponseWriter, r *http.Request) error
}

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (userH UserHandler) Register(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var req models.RegisterInput

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return httpx.BadRequest("invalid JSON")
	}

	user, err := userH.UserService.Register(ctx, req)
	if err != nil {
		log.Printf("failed to add user: %v", err)
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(w, http.StatusCreated, user)
	return nil
}

func (userH *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	users, err := userH.UserService.GetAllUsers(ctx)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(w, http.StatusOK, users)
	return nil
}

func (userH *UserHandler) Login(w http.ResponseWriter, r *http.Request) error {
	return nil
}
