package repository

import (
	"enactus/internal/models"
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
