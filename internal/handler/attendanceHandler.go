package handler

import (
	"enactus/internal/helpers"
	"enactus/internal/httpx"
	"enactus/internal/models/inputs"
	"enactus/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type AttendanceHandlerInterface interface {
	AddAttendance(w http.ResponseWriter, r *http.Request) error
	GetAllAttendances(w http.ResponseWriter, r *http.Request) error
	UpdateAttendance(w http.ResponseWriter, r *http.Request) error
	DeleteAttendance(w http.ResponseWriter, r *http.Request) error
}

type AttendanceHandler struct {
	AttendanceS *service.AttendanceService
}

func NewAttendanceHandler(attendanceS *service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{AttendanceS: attendanceS}
}

func (attendanceH *AttendanceHandler) AddAttendance(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req inputs.AddAttendanceInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return httpx.BadRequest("invalid request body")
	}

	userId, ok := r.Context().Value("user_id").(int)
	if !ok {
		return httpx.BadRequest("invalid user id")
	}

	req.Creator = userId
	added, err := attendanceH.AttendanceS.AddAttendance(ctx, &req)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(w, http.StatusCreated, added)
	return nil
}

func (attendanceH *AttendanceHandler) GetAllAttendances(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	attendances, err := attendanceH.AttendanceS.GetAllAttendances(ctx)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(w, http.StatusOK, attendances)
	return nil
}

func (attendanceH *AttendanceHandler) UpdateAttendance(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return httpx.BadRequest("invalid attendance id")
	}

	var req inputs.UpdateAttendanceInput
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return httpx.BadRequest("invalid request body")
	}

	userId, ok := r.Context().Value("user_id").(int)
	if !ok {
		return httpx.BadRequest("invalid user id")
	}

	req.MarkedBy = &userId

	v := helpers.NewValidator()
	errors := helpers.Validate(req, v)
	if errors != nil {
		return httpx.BadRequest("invalid status")
	}

	err = attendanceH.AttendanceS.UpdateAttendance(ctx, id, &req)
	if err != nil {
		return httpx.InternalError(err)
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (attendanceH *AttendanceHandler) DeleteAttendance(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return httpx.BadRequest("invalid id")
	}

	err = attendanceH.AttendanceS.DeleteAttendance(ctx, id)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(w, http.StatusNoContent, nil)
	return nil
}
