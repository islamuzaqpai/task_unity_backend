package service

import (
	"context"
	"enactus/internal/apperrors"
	"enactus/internal/models"
	"enactus/internal/models/inputs"
	"enactus/internal/repository"
	"errors"
	"fmt"
)

type TaskServiceInterface interface {
	AddTask(ctx context.Context, task *models.Task) (*models.Task, error)
	GetAllTasks(ctx context.Context) ([]models.Task, error)
	GetTaskById(ctx context.Context, id int) (*models.Task, error)
	GetAllTasksByAssigneeId(ctx context.Context, creatorId int) ([]models.Task, error)
	UpdateTask(ctx context.Context, userId, taskId int, role string, in inputs.UpdateTaskInput) (*models.Task, error)
	DeleteTask(ctx context.Context, taskId, userId int, role string) error
}

type TaskService struct {
	TaskRepo *repository.TaskRepository
	UserRepo *repository.UserRepository
}

func NewTaskService(taskR *repository.TaskRepository, userR *repository.UserRepository) *TaskService {
	return &TaskService{TaskRepo: taskR, UserRepo: userR}
}

func (taskS *TaskService) AddTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	creator, err := taskS.UserRepo.GetUserById(ctx, task.CreatorId)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrCreatorNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	assignee, err := taskS.UserRepo.GetUserById(ctx, task.AssigneeId)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrAssigneeNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if creator.DepartmentId == nil || assignee.DepartmentId == nil {
		return nil, fmt.Errorf("creator and assignee must be in the same department")
	}

	if *creator.DepartmentId != *assignee.DepartmentId {
		return nil, fmt.Errorf("creator and assignee must be in the same department")
	}

	task.DepartmentId = *creator.DepartmentId

	addedTask, err := taskS.TaskRepo.AddTask(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("failed to add a task: %w", err)
	}

	return addedTask, nil
}

func (taskS *TaskService) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	tasks, err := taskS.TaskRepo.GetAllTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all tasks: %w", err)
	}

	return tasks, nil
}

func (taskS *TaskService) GetAllTasksByAssigneeId(ctx context.Context, assigneeId int) ([]models.Task, error) {
	tasks, err := taskS.TaskRepo.GetAllTasksByAssigneeId(ctx, assigneeId)
	if err != nil {
		return nil, fmt.Errorf("failed to get all tasks: %w", err)
	}

	return tasks, err
}

func (taskS *TaskService) GetTaskById(ctx context.Context, id int) (*models.Task, error) {
	task, err := taskS.TaskRepo.GetTaskById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get a task: %w", err)
	}

	return task, nil
}

func (taskS *TaskService) UpdateTask(ctx context.Context, userId, taskId int, role string, in inputs.UpdateTaskInput) (*models.Task, error) {
	task, err := taskS.GetTaskById(ctx, taskId)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	if role != "admin" && role != "manager" && task.CreatorId != userId {
		return nil, apperrors.ErrForbidden
	}

	updated, err := taskS.TaskRepo.UpdateTask(ctx, taskId, in)
	if err != nil {
		return nil, fmt.Errorf("failed to update a task: %w", err)
	}

	return updated, nil
}

func (taskS *TaskService) DeleteTask(ctx context.Context, taskId, userId int, role string) error {
	task, err := taskS.GetTaskById(ctx, taskId)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.ErrNotFound
		}
		return fmt.Errorf("failed to get task: %w", err)
	}

	if role != "admin" && role != "manager" && task.CreatorId != userId {
		return apperrors.ErrForbidden
	}

	err = taskS.TaskRepo.DeleteTask(ctx, taskId)
	if err != nil {
		return fmt.Errorf("failed to delete a task: %w", err)
	}

	return nil
}
