package repository

import (
	"context"
	"enactus/internal/models"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepositoryInterface interface {
	AddTask(ctx context.Context, task *models.Task) error
	GetTaskById(ctx context.Context, id int) (*models.Task, error)
	GetAllTasks(ctx context.Context) ([]models.Task, error)
	UpdateTask(ctx context.Context, id int, newTask models.Task) error
	DeleteTask(ctx context.Context, id int) error
}

type TaskRepository struct {
	Pool *pgxpool.Pool
}

func (taskRepo *TaskRepository) AddTask(ctx context.Context, task *models.Task) error {

	row := taskRepo.Pool.QueryRow(ctx,
		"INSERT INTO tasks (title, description, deadline, department_id, creator_id, assignee_id, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, title, description, deadline, department_id, creator_id, assignee_id, status",
		task.Title,
		task.Description,
		task.Deadline,
		task.DepartmentId,
		task.CreatorId,
		task.AssigneeId,
		task.Status,
	)

	err := row.Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.Deadline,
		&task.DepartmentId, &task.CreatorId,
		&task.AssigneeId,
		&task.Status,
	)

	if err != nil {
		return fmt.Errorf("failed to add a task: %w", err)
	}

	return nil
}

func (taskRepo *TaskRepository) GetTaskById(ctx context.Context, id int) (*models.Task, error) {
	row := taskRepo.Pool.QueryRow(ctx,
		"SELECT id, title, description, deadline, department_id, creator_id, assignee_id, status, created_at, updated_at FROM tasks WHERE id = $1 AND deleted_at IS NULL ",
		id,
	)

	var task models.Task

	err := row.Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.Deadline,
		&task.DepartmentId,
		&task.CreatorId,
		&task.AssigneeId,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return &task, nil
}

func (taskRepo *TaskRepository) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	rows, err := taskRepo.Pool.Query(ctx,
		"SELECT id, title, description, deadline, department_id, creator_id, assignee_id, status, created_at, updated_at, deleted_at FROM tasks WHERE deleted_at is null")

	if err != nil {
		return nil, fmt.Errorf("failed to select all tasks: %w", err)
	}

	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task

		err = rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.Deadline,
			&task.DepartmentId,
			&task.CreatorId,
			&task.AssigneeId,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.DeletedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		tasks = append(tasks, task)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return tasks, nil
}

func (taskRepo *TaskRepository) UpdateTask(ctx context.Context, id int, newTask models.Task) error {
	row := taskRepo.Pool.QueryRow(ctx,
		"UPDATE tasks SET title = $1, description = $2, deadline = $3, assignee_id = $4, status = $5, updated_at = now() WHERE id = $6 AND deleted_at IS NULL RETURNING id",
		newTask.Title,
		newTask.Description,
		newTask.Deadline,
		newTask.AssigneeId,
		newTask.Status,
		id,
	)

	var returnedId int
	err := row.Scan(
		&returnedId,
	)

	if err != nil {
		return fmt.Errorf("failed to scan: %w", err)
	}

	return nil
}

func (taskRepo *TaskRepository) DeleteTask(ctx context.Context, id int) error {
	_, err := taskRepo.Pool.Exec(ctx,
		"UPDATE tasks SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL",
		id)

	if err != nil {
		return fmt.Errorf("failed to delete a task: %w", err)
	}

	return nil
}
