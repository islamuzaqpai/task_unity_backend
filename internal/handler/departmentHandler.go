package handler

import (
	"enactus/internal/httpx"
	"enactus/internal/models"
	"enactus/internal/service"
	"encoding/json"
	"net/http"
)

type DepartmentHandlerInterface interface {
	AddDepartment(w http.ResponseWriter, r *http.Request) error
}

type DepartmentHandler struct {
	DepartmentS *service.DepartmentService
}

func NewDepartmentHandler(departmentS *service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{DepartmentS: departmentS}
}

func (departmentH *DepartmentHandler) AddDepartment(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var department models.Department
	err := json.NewDecoder(r.Body).Decode(&department)
	if err != nil {
		return httpx.BadRequest("invalid request body")
	}

	added, err := departmentH.DepartmentS.AddDepartment(ctx, &department)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(w, http.StatusOK, added)
	return nil
}
