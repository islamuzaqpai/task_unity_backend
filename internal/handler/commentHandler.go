package handler

import (
	"enactus/internal/apperrors"
	"enactus/internal/httpx"
	"enactus/internal/models/inputs"
	"enactus/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type CommentHandlerInterface interface {
	AddComment(w http.ResponseWriter, r *http.Request) error
	GetAllComments(w http.ResponseWriter, r *http.Request) error
	UpdateComment(w http.ResponseWriter, r *http.Request) error
	DeleteComment(w http.ResponseWriter, r *http.Request) error
}

type CommentHandler struct {
	CommentS *service.CommentService
}

func NewCommentHandler(commentS *service.CommentService) *CommentHandler {
	return &CommentHandler{CommentS: commentS}
}

func (commentH *CommentHandler) AddComment(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	var req inputs.AddCommentInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return httpx.BadRequest("invalid request body")
	}

	creatorId, ok := r.Context().Value("user_id").(int)
	if !ok {
		return httpx.BadRequest("invalid user id")
	}

	req.CreatorId = creatorId

	added, err := commentH.CommentS.AddComment(ctx, &req)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(w, http.StatusCreated, added)
	return nil
}

func (commentH *CommentHandler) GetAllComments(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	comments, err := commentH.CommentS.GetAllComments(ctx)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(w, http.StatusOK, comments)
	return nil
}

func (commentH *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return httpx.BadRequest("invalid ID")
	}

	userId, ok := r.Context().Value("user_id").(int)
	if !ok {
		return httpx.Unauthorized("user_id missing")
	}

	var req inputs.UpdateCommentInput
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return httpx.BadRequest("invalid request body")
	}

	req.UserId = userId

	err = commentH.CommentS.UpdateComment(ctx, id, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return httpx.NotFound("comment")
		}

		return httpx.InternalError(err)
	}

	httpx.WriteJSON(w, http.StatusOK, "ok")
	return nil
}

func (commentH *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return httpx.BadRequest("invalid ID")
	}

	userId, ok := r.Context().Value("user_id").(int)
	if !ok {
		return httpx.BadRequest("user_id missing")
	}

	err = commentH.CommentS.DeleteComment(ctx, userId, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return httpx.NotFound("comment")
		}

		return httpx.InternalError(err)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
