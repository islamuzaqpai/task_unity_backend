package repository

import (
	"context"
	"database/sql"
	"enactus/internal/apperrors"
	"enactus/internal/models"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentRepositoryInterface interface {
	AddComment(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	GetCommentById(ctx context.Context, id int) (*models.Comment, error)
	GetAllComments(ctx context.Context) ([]models.Comment, error)
	UpdateComment(ctx context.Context, id int, newComment models.Comment) error
	DeleteComment(ctx context.Context, id int) error
}

type CommentRepository struct {
	Pool *pgxpool.Pool
}

func NewCommentRepository(pool *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{Pool: pool}
}

func (commentRepo *CommentRepository) AddComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	query := `INSERT INTO tasks_comments (comment, task_id, creator_id) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`

	err := commentRepo.Pool.QueryRow(ctx, query,
		comment.Comment,
		comment.TaskId,
		comment.CreatorId,
	).Scan(
		&comment.Id,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return comment, nil
}

func (commentRepo *CommentRepository) GetCommentById(ctx context.Context, id int) (*models.Comment, error) {
	query := `SELECT id, comment, task_id, creator_id, created_at, updated_at, deleted_at FROM tasks_comments WHERE id = $1 AND deleted_at IS NULL`

	var comment models.Comment
	err := commentRepo.Pool.QueryRow(ctx, query, id).Scan(
		&comment.Id,
		&comment.Comment,
		&comment.TaskId,
		&comment.CreatorId,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return &comment, nil
}

func (commentRepo *CommentRepository) GetAllComments(ctx context.Context) ([]models.Comment, error) {
	query := `SELECT id, comment, task_id, creator_id, created_at, updated_at FROM tasks_comments WHERE deleted_at IS NULL`

	rows, err := commentRepo.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all comments: %w", err)
	}
	defer rows.Close()

	var comments []models.Comment

	for rows.Next() {
		var comment models.Comment

		err = rows.Scan(
			&comment.Id,
			&comment.Comment,
			&comment.TaskId,
			&comment.CreatorId,
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

func (commentRepo *CommentRepository) UpdateComment(ctx context.Context, id int, comment string) error {
	query := `UPDATE tasks_comments SET comment = $1, updated_at = now() WHERE id = $2 AND deleted_at IS NULL RETURNING id`

	var returnedId int
	err := commentRepo.Pool.QueryRow(ctx, query, comment, id).Scan(
		&returnedId,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperrors.ErrNotFound
		}
		return fmt.Errorf("failed to scan: %w", err)
	}

	return nil
}

func (commentRepo *CommentRepository) DeleteComment(ctx context.Context, id int) error {
	query := `UPDATE tasks_comments SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`

	res, err := commentRepo.Pool.Exec(ctx, query, id)

	if res.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}

	if err != nil {
		return fmt.Errorf("failed to delete a comment: %w", err)
	}

	return nil
}
