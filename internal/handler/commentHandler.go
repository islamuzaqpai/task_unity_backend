package handler

import (
	"enactus/internal/httpx"
	"enactus/internal/models/inputs"
	"enactus/internal/service"
	"encoding/json"
	"net/http"
)

type CommentHandlerInterface interface {
	AddComment(w http.ResponseWriter, r *http.Request) error
	GetAllComments(w http.ResponseWriter, r *http.Request) error
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
