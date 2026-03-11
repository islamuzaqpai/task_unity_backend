package httpx

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Status  int               `json:"-"`
	Details map[string]string `json:"details,omitempty"`
	Err     error             `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}

	return e.Message
}

func NotFound(resource string) *AppError {
	return &AppError{
		Code:    "NOT_FOUND",
		Message: fmt.Sprintf("%s not found", resource),
		Status:  http.StatusNotFound,
	}
}

func BadRequest(message string) *AppError {
	return &AppError{
		Code:    "BAD_REQUEST",
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

func ValidationError(err error) *AppError {
	return &AppError{
		Code:    "INTERNAL_ERROR",
		Message: "internal server error",
		Status:  http.StatusUnprocessableEntity,
		Err:     err,
	}
}

func InternalError(err error) *AppError {
	return &AppError{
		Code:    "INTERNAL_ERROR",
		Message: "internal server error",
		Status:  http.StatusInternalServerError,
		Err:     err,
	}
}

type AppHandler func(w http.ResponseWriter, r *http.Request) error

func WrapHandler(h AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err == nil {
			return
		}

		var appErr *AppError
		if errors.As(err, &appErr) {
			WriteJSON(w, appErr.Status, appErr)
			return
		}

		internalErr := InternalError(err)
		WriteJSON(w, internalErr.Status, appErr)
	}
}
