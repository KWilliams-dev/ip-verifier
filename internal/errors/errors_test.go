package errors

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name        string
		appErr      *AppError
		expectedMsg string
	}{
		{
			name: "with internal error",
			appErr: &AppError{
				Code:        http.StatusBadRequest,
				Message:     "Invalid input",
				InternalErr: fmt.Errorf("field is required"),
			},
			expectedMsg: "Invalid input: field is required",
		},
		{
			name: "without internal error",
			appErr: &AppError{
				Code:        http.StatusNotFound,
				Message:     "Resource not found",
				InternalErr: nil,
			},
			expectedMsg: "Resource not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedMsg, tt.appErr.Error())
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	internalErr := fmt.Errorf("internal error")
	appErr := &AppError{
		Code:        http.StatusInternalServerError,
		Message:     "Server error",
		InternalErr: internalErr,
	}

	assert.Equal(t, internalErr, appErr.Unwrap())
}

func TestNewValidationError(t *testing.T) {
	internalErr := fmt.Errorf("validation failed")
	err := NewValidationError("Invalid request", internalErr)

	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Equal(t, "Invalid request", err.Message)
	assert.Equal(t, internalErr, err.InternalErr)
}

func TestNewNotFoundError(t *testing.T) {
	internalErr := fmt.Errorf("record not found")
	err := NewNotFoundError("Resource not found", internalErr)

	assert.Equal(t, http.StatusNotFound, err.Code)
	assert.Equal(t, "Resource not found", err.Message)
	assert.Equal(t, internalErr, err.InternalErr)
}

func TestNewInternalError(t *testing.T) {
	internalErr := fmt.Errorf("database connection failed")
	err := NewInternalError("Internal server error", internalErr)

	assert.Equal(t, http.StatusInternalServerError, err.Code)
	assert.Equal(t, "Internal server error", err.Message)
	assert.Equal(t, internalErr, err.InternalErr)
}

func TestGetHTTPStatus(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "validation error",
			err:            NewValidationError("Invalid input", nil),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "not found error",
			err:            NewNotFoundError("Not found", nil),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "internal error",
			err:            NewInternalError("Server error", nil),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "generic error",
			err:            fmt.Errorf("generic error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := GetHTTPStatus(tt.err)
			assert.Equal(t, tt.expectedStatus, status)
		})
	}
}

func TestGetMessage(t *testing.T) {
	tests := []struct {
		name            string
		err             error
		expectedMessage string
	}{
		{
			name:            "app error",
			err:             NewValidationError("Invalid input", nil),
			expectedMessage: "Invalid input",
		},
		{
			name:            "generic error",
			err:             fmt.Errorf("generic error"),
			expectedMessage: "An internal error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := GetMessage(tt.err)
			assert.Equal(t, tt.expectedMessage, message)
		})
	}
}

func TestIsValidationError(t *testing.T) {
	assert.True(t, IsValidationError(NewValidationError("test", nil)))
	assert.False(t, IsValidationError(NewNotFoundError("test", nil)))
	assert.False(t, IsValidationError(NewInternalError("test", nil)))
	assert.False(t, IsValidationError(errors.New("generic")))
}

func TestIsNotFoundError(t *testing.T) {
	assert.True(t, IsNotFoundError(NewNotFoundError("test", nil)))
	assert.False(t, IsNotFoundError(NewValidationError("test", nil)))
	assert.False(t, IsNotFoundError(NewInternalError("test", nil)))
	assert.False(t, IsNotFoundError(errors.New("generic")))
}

func TestIsInternalError(t *testing.T) {
	assert.True(t, IsInternalError(NewInternalError("test", nil)))
	assert.False(t, IsInternalError(NewValidationError("test", nil)))
	assert.False(t, IsInternalError(NewNotFoundError("test", nil)))
	assert.False(t, IsInternalError(errors.New("generic")))
}
