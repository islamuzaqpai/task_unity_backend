package repository

import (
	"context"
	"enactus/internal/models"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentRepositoryInterface interface {
	AddComment(ctx context.Context, comment *models.Comment) error
	GetCommentById(ctx context.Context, id int) (*models.Comment, error)
	GetAllComments(ctx context.Context) ([]models.Comment, error)
	UpdateComment(ctx context.Context, id int, newComment models.Comment) error
	DeleteComment(ctx context.Context, id int) error
}

type CommentRepository struct {
	Pool *pgxpool.Pool
}

func (commentRepo *CommentRepository) AddComment(ctx context.Context, comment *models.Comment) error {
	row := commentRepo.Pool.QueryRow(ctx,
		"INSERT INTO comments (description, task_id, user_id) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at",
		comment.Description,
		comment.TaskId,
		comment.UserId,
	)

	err := row.Scan(
		&comment.Id,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to scan: %w", err)
	}

	return nil
}

func (commentRepo *CommentRepository) GetCommentById(ctx context.Context, id int) (*models.Comment, error) {
	row := commentRepo.Pool.QueryRow(ctx,
		"SELECT id, description, task_id, user_id, created_at, updated_at, deleted_at FROM comments WHERE id = $1 AND deleted_at IS NULL",
		id,
	)

	var comment models.Comment
	err := row.Scan(
		&comment.Id,
		&comment.Description,
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

func (commentRepo *CommentRepository) GetAllComments(ctx context.Context) ([]models.Comment, error) {
	rows, err := commentRepo.Pool.Query(ctx,
		"SELECT id, description, task_id, user_id, created_at, updated_at FROM comments WHERE deleted_at IS NULL",
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get all comments: %w", err)
	}
	defer rows.Close()

	var comments []models.Comment

	for rows.Next() {
		var comment models.Comment

		err = rows.Scan(
			&comment.Id,
			&comment.Description,
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
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return comments, nil
}

func (commentRepo *CommentRepository) UpdateComment(ctx context.Context, id int, newComment models.Comment) error {
	row := commentRepo.Pool.QueryRow(ctx,
		"UPDATE comments SET description = $1, task_id = $2, user_id = $3, updated_at = now() WHERE id = $4 AND deleted_at IS NULL RETURNING id",
		newComment.Description,
		newComment.TaskId,
		newComment.UserId,
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

func (commentRepo *CommentRepository) DeleteComment(ctx context.Context, id int) error {
	_, err := commentRepo.Pool.Exec(ctx,
		"UPDATE comments SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL",
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to delete a comment: %w", err)
	}

	return nil
}
