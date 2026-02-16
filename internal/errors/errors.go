package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError represents a custom application error with HTTP status code
type AppError struct {
	Code        int    // HTTP status code
	Message     string // User-facing message
	InternalErr error  // Internal error for logging
}

func (e *AppError) Error() string {
	if e.InternalErr != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.InternalErr)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.InternalErr
}

// NewValidationError creates a validation error (400 Bad Request)
func NewValidationError(message string, err error) *AppError {
	return &AppError{
		Code:        http.StatusBadRequest,
		Message:     message,
		InternalErr: err,
	}
}

// NewNotFoundError creates a not found error (404 Not Found)
func NewNotFoundError(message string, err error) *AppError {
	return &AppError{
		Code:        http.StatusNotFound,
		Message:     message,
		InternalErr: err,
	}
}

// NewInternalError creates an internal server error (500 Internal Server Error)
func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Code:        http.StatusInternalServerError,
		Message:     message,
		InternalErr: err,
	}
}

// GetHTTPStatus extracts HTTP status code from error
func GetHTTPStatus(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return http.StatusInternalServerError
}

// GetMessage extracts user-facing message from error
func GetMessage(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Message
	}
	return "An internal error occurred"
}

// IsValidationError checks if error is a validation error
func IsValidationError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == http.StatusBadRequest
	}
	return false
}

// IsNotFoundError checks if error is a not found error
func IsNotFoundError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == http.StatusNotFound
	}
	return false
}

// IsInternalError checks if error is an internal error
func IsInternalError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == http.StatusInternalServerError
	}
	return false
}
