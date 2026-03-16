package handler

import (
	"enactus/internal/httpx"
	"enactus/internal/models"
	"enactus/internal/models/inputs"
	"enactus/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type TaskHandlerInterface interface {
	AddTask(w http.ResponseWriter, r *http.Request) error
	GetAllTasks(w http.ResponseWriter, r *http.Request) error
	GetAllTasksByAssigneeId(w http.ResponseWriter, r *http.Request) error
	GetTaskById(w http.ResponseWriter, r *http.Request) error
	UpdateTask(w http.ResponseWriter, r *http.Request) error
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

func (taskH *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	tasks, err := taskH.TaskS.GetAllTasks(ctx)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(w, http.StatusOK, tasks)
	return nil
}

func (taskH *TaskHandler) GetAllTasksByAssigneeId(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	userId := r.Context().Value("user_id").(int)
	tasks, err := taskH.TaskS.GetAllTasksByAssigneeId(ctx, userId)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(w, http.StatusOK, tasks)
	return nil
}

func (taskH *TaskHandler) GetTaskById(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return httpx.BadRequest("invalid task ID")
	}

	task, err := taskH.TaskS.GetTaskById(ctx, id)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(w, http.StatusOK, task)
	return nil
}

func (taskH *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	userId := r.Context().Value("user_id").(int)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return httpx.BadRequest("invalid task ID")
	}

	var req inputs.UpdateTaskInput
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return httpx.BadRequest("invalid request body")
	}

	updated, err := taskH.TaskS.UpdateTask(ctx, userId, id, req)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(w, http.StatusOK, updated)
	return nil
}
