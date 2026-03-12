package handler

import (
	"enactus/internal/httpx"
	"enactus/internal/models"
	"enactus/internal/models/inputs"
	"enactus/internal/service"
	"encoding/json"
	"net/http"
)

type TaskHandlerInterface interface {
	AddTask(w http.ResponseWriter, r *http.Request) error
}

type TaskHandler struct {
	TaskS *service.TaskService
}

func NewTaskHandler(taskS *service.TaskService) *TaskHandler {
	return &TaskHandler{TaskS: taskS}
}

func (taskH *TaskHandler) AddTask(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req inputs.AddTaskInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return httpx.BadRequest("invalid request body")
	}

	userId := r.Context().Value("user_id").(int)
	task := models.Task{
		Title:        req.Title,
		Description:  req.Description,
		Deadline:     req.Deadline,
		DepartmentId: req.DepartmentId,
		CreatorId:    userId,
		AssigneeId:   req.AssigneeId,
		Status:       req.Status,
	}

	addedTask, err := taskH.TaskS.AddTask(ctx, &task)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(w, http.StatusOK, addedTask)
	return nil
}
