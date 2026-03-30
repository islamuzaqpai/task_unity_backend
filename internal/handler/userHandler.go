package handler

import (
	"enactus/internal/apperrors"
	"enactus/internal/httpx"
	"enactus/internal/models/inputs"
	"enactus/internal/service"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type UserHandlerInterface interface {
	Register(w http.ResponseWriter, r *http.Request) error
	GetAllUsers(w http.ResponseWriter, r *http.Request) error
	GetUserById(w http.ResponseWriter, r *http.Request) error
	Login(w http.ResponseWriter, r *http.Request) error
	UpdateUserProfile(w http.ResponseWriter, r *http.Request) error
	UpdateUserPassword(w http.ResponseWriter, r *http.Request) error
	DeleteUser(w http.ResponseWriter, r *http.Request) error
}

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (userH UserHandler) Register(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var req inputs.RegisterInput

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return httpx.BadRequest("invalid JSON")
	}

	user, err := userH.UserService.Register(ctx, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrEmailAlreadyExists) {
			return httpx.UserAlreadyExists()
		}

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
	ctx := r.Context()

	var req inputs.LoginInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return httpx.BadRequest("invalid request body")
	}

	tokenStr, err := userH.UserService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return httpx.Unauthorized("invalid email or password")
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"token": tokenStr})
	return nil
}

func (userH *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return httpx.BadRequest("invalid ID")
	}

	user, err := userH.UserService.GetUserById(ctx, id)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(w, http.StatusOK, user)
	return nil
}

func (userH *UserHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return httpx.BadRequest("invalid ID")
	}

	var req inputs.UpdateUserProfileInput
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return httpx.BadRequest("invalid request body")
	}

	updated, err := userH.UserService.UpdateUserProfile(ctx, id, req)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(w, http.StatusOK, updated)
	return nil
}

func (userH *UserHandler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return httpx.BadRequest("invalid ID")
	}

	var req inputs.UpdatePasswordInput
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return httpx.BadRequest("invalid request body")
	}

	err = userH.UserService.UpdateUserPassword(ctx, id, req.Password)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(w, http.StatusOK, "OK")
	return nil
}

func (userH *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return httpx.BadRequest("invalid ID")
	}

	err = userH.UserService.DeleteUser(ctx, id)
	if err != nil {
		return httpx.InternalError(err)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
