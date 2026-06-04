package apperror

import (
	"fmt"
	"net/http"
)

// AppError represents a typed application error that maps cleanly to HTTP.
type AppError struct {
	Code       string         `json:"code"`
	Message    string         `json:"message"`
	HTTPStatus int            `json:"-"`
	Details    map[string]any `json:"details"`
	err        error          `json:"-"` // wrapped error for internal logging
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap allows errors.Is and errors.As to work.
func (e *AppError) Unwrap() error {
	return e.err
}

// Constructors

func New(httpStatus int, code, message string, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Details:    make(map[string]any),
		err:        err,
	}
}

func NotFound(message string, err error) *AppError {
	return New(http.StatusNotFound, "RESOURCE_NOT_FOUND", message, err)
}

func Unauthorized(message string, err error) *AppError {
	return New(http.StatusUnauthorized, "UNAUTHORIZED", message, err)
}

func Forbidden(message string, err error) *AppError {
	return New(http.StatusForbidden, "FORBIDDEN", message, err)
}

func BadRequest(message string, err error) *AppError {
	return New(http.StatusBadRequest, "BAD_REQUEST", message, err)
}

func Conflict(message string, err error) *AppError {
	return New(http.StatusConflict, "CONFLICT", message, err)
}

func Internal(err error) *AppError {
	return New(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "An internal server error occurred", err)
}

func TooManyRequests(message string, err error) *AppError {
	return New(http.StatusTooManyRequests, "TOO_MANY_REQUESTS", message, err)
}

func Validation(details map[string]any, err error) *AppError {
	e := New(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err)
	if details != nil {
		e.Details = details
	}
	return e
}
