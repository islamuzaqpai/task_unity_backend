package service

import (
	"context"
	"enactus/internal/models"
	"enactus/internal/repository"
	"fmt"
)

type CommentServiceInterface interface {
	AddComment(ctx context.Context, comment *models.Comment) (*models.Comment, error)
}

type CommentService struct {
	CommentRepo *repository.CommentRepository
}

func (commentS *CommentService) AddComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	err := commentS.CommentRepo.AddComment(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("failed to add a comment: %w", err)
	}

	return comment, nil
}
