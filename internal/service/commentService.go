package service

import (
	"context"
	"enactus/internal/models"
	"enactus/internal/repository"
	"fmt"
)

type CommentServiceInterface interface {
	AddComment(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	GetAllComments(ctx context.Context) ([]models.Comment, error)
	GetCommentById(ctx context.Context, id int) (*models.Comment, error)
	UpdateComment(ctx context.Context, id int, description string) error
	DeleteComment(ctx context.Context, id int) error
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

func (commentS *CommentService) GetAllComments(ctx context.Context) ([]models.Comment, error) {
	comments, err := commentS.CommentRepo.GetAllComments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all comments: %w", err)
	}

	return comments, nil
}

func (commentS *CommentService) GetCommentById(ctx context.Context, id int) (*models.Comment, error) {
	comment, err := commentS.CommentRepo.GetCommentById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get a comment: %w", err)
	}

	return comment, err
}

func (commentS *CommentService) UpdateComment(ctx context.Context, id int, description string) error {
	err := commentS.CommentRepo.UpdateComment(ctx, id, description)
	if err != nil {
		return fmt.Errorf("failed to update a comment: %w", err)
	}

	return nil
}

func (commentS *CommentService) DeleteComment(ctx context.Context, id int) error {
	err := commentS.CommentRepo.DeleteComment(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete a comment: %w", err)
	}

	return nil
}
