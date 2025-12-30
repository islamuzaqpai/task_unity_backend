package repository

import (
	"context"
	"enactus/internal/models"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentRepositoryInterface interface {
	AddComment(comment *models.TaskComment) (*models.TaskComment, error)
	GetCommentById(id int) (*models.TaskComment, error)
	GetAllComments() ([]models.TaskComment, error)
	UpdateComment(id int, newComment models.TaskComment) (*models.TaskComment, error)
	DeleteComment(id int) error
}

type CommentRepository struct {
	Pool *pgxpool.Pool
}

func (commentRepo *CommentRepository) AddComment(comment *models.TaskComment) (*models.TaskComment, error) {
	row := commentRepo.Pool.QueryRow(context.Background(),
		"INSERT INTO tasks_comments (comment, task_id, user_id) VALUES ($1, $2, $3) RETURNING id, comment, task_id, user_id, created_at, updated_at",
		comment.Comment,
		comment.TaskId,
		comment.UserId,
	)

	err := row.Scan(
		&comment.Id,
		&comment.Comment,
		&comment.TaskId,
		&comment.UserId,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return comment, nil
}

func (commentRepo *CommentRepository) GetCommentById(id int) (*models.TaskComment, error) {
	row := commentRepo.Pool.QueryRow(context.Background(),
		"SELECT id, comment, task_id, user_id, created_at, updated_at, deleted_at FROM tasks_comments WHERE id = $1",
		id,
	)

	var comment models.TaskComment
	err := row.Scan(
		&comment.Id,
		&comment.Comment,
		&comment.TaskId,
		&comment.UserId,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.DeletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return &comment, nil
}

func (commentRepo *CommentRepository) GetAllComments() ([]models.TaskComment, error) {
	rows, err := commentRepo.Pool.Query(context.Background(),
		"SELECT id, comment, task_id, user_id, created_at, updated_at FROM tasks_comments WHERE deleted_at IS NULL",
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get all comments: %w", err)
	}

	var comments []models.TaskComment

	for rows.Next() {
		var comment models.TaskComment

		err = rows.Scan(
			&comment.Id,
			&comment.Comment,
			&comment.TaskId,
			&comment.UserId,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		comments = append(comments, comment)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iteration: %w", err)
	}

	defer rows.Close()

	return comments, nil
}

func (commentRepo *CommentRepository) UpdateComment(id int, newComment models.TaskComment) (*models.TaskComment, error) {
	return nil, nil
}

func (commentRepo *CommentRepository) DeleteComment(id int) error {
	return nil
}
