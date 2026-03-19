package repository

import (
	"context"
	"database/sql"
	"enactus/internal/apperrors"
	"enactus/internal/models"
	"enactus/internal/models/inputs"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepositoryInterface interface {
	AddTask(ctx context.Context, task *models.Task) (*models.Task, error)
	GetTaskById(ctx context.Context, id int) (*models.Task, error)
	GetAllTasks(ctx context.Context) ([]models.Task, error)
	GetAllTasksByAssigneeId(ctx context.Context, assigneeId int) ([]models.Task, error)
	UpdateTask(ctx context.Context, id int, in inputs.UpdateTaskInput) (*models.Task, error)
	DeleteTask(ctx context.Context, id int) error
}
type TaskRepository struct {
	Pool *pgxpool.Pool
}

func NewTaskRepo(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{Pool: pool}
}

func (taskRepo *TaskRepository) AddTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	query := `INSERT INTO tasks (title, description, deadline, department_id, creator_id, assignee_id, status) VALUES ($1, $2, $3, $4, $5, $6, $7)
				RETURNING id, created_at, updated_at, deleted_at`

	err := taskRepo.Pool.QueryRow(ctx, query,
		task.Title,
		task.Description,
		task.Deadline,
		task.DepartmentId,
		task.CreatorId,
		task.AssigneeId,
		task.Status,
	).Scan(
		&task.Id,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.DeletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to add a task: %w", err)
	}

	return task, nil
}

func (taskRepo *TaskRepository) GetTaskById(ctx context.Context, id int) (*models.Task, error) {
	query := `SELECT id, title, description, deadline, department_id, creator_id, assignee_id, status, created_at, updated_at FROM tasks WHERE id = $1 AND deleted_at IS NULL`
	var task models.Task
	err := taskRepo.Pool.QueryRow(ctx, query, id).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}

		return nil, fmt.Errorf("failed to get task by id: %w", err)
	}

	return &task, nil
}

func (taskRepo *TaskRepository) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	query := `SELECT id, title, description, deadline, department_id, creator_id, assignee_id, status, created_at, updated_at, deleted_at FROM tasks WHERE deleted_at is null`

	rows, err := taskRepo.Pool.Query(ctx, query)

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

func (taskRepo *TaskRepository) GetAllTasksByAssigneeId(ctx context.Context, assigneeId int) ([]models.Task, error) {
	query := `SELECT id, title, description, deadline, department_id, creator_id, assignee_id, status, created_at, updated_at, deleted_at FROM tasks WHERE deleted_at is null AND assignee_id = $1`

	rows, err := taskRepo.Pool.Query(ctx, query, assigneeId)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
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

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration error: %w", rows.Err())
	}

	return tasks, nil
}

func (taskRepo *TaskRepository) UpdateTask(ctx context.Context, id int, in inputs.UpdateTaskInput) (*models.Task, error) {
	query := `UPDATE tasks SET title = COALESCE($1, title),
								description = COALESCE($2, description),
								deadline = COALESCE($3, deadline),
								assignee_id = COALESCE($4, assignee_id),
								status = COALESCE($5::task_status, status),
								updated_at = NOW()
				WHERE id = $6 AND deleted_at IS NULL
				RETURNING id, title, description, deadline, department_id, creator_id, assignee_id, status, created_at, updated_at, deleted_at`

	var task models.Task
	err := taskRepo.Pool.QueryRow(ctx, query,
		in.Title,
		in.Description,
		in.Deadline,
		in.AssigneeId,
		in.Status,
		id,
	).Scan(
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

	return &task, nil
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
