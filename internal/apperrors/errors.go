package apperrors

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrWeakPassword       = errors.New("weak password")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrAssigneeNotFound   = errors.New("assignee user not found")
	ErrCreatorNotFound    = errors.New("creator user not found")
)
