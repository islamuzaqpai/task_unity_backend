package service

import (
	"context"
	"enactus/internal/models"
	"enactus/internal/repository"
	"fmt"
)

type TaskServiceInterface interface {
	AddTask(ctx context.Context, task *models.Task) (*models.Task, error)
	GetAllTasks(ctx context.Context) ([]models.Task, error)
}

type TaskService struct {
	TaskRepo *repository.TaskRepository
}

func (taskS *TaskService) AddTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	err := taskS.TaskRepo.AddTask(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("failed to add a task: %w", err)
	}

	return task, nil
}

func (taskS *TaskService) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	tasks, err := taskS.TaskRepo.GetAllTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all tasks: %w", err)
	}

	return tasks, nil
}
