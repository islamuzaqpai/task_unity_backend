package apperrors

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrWeakPassword       = errors.New("weak password")
	ErrUnauthorized       = errors.New("unauthorized")
)
