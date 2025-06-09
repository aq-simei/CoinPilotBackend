package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/aq-simei/coin-pilot/internal/config/logger"
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
	Status  int    `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return "[" + e.Code + "] " + e.Message + ": " + e.Err.Error()
	}
	return "[" + e.Code + "] " + e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// ========== Constructors ==========

func New(code string, message string, status int) *AppError {
	logger.Warn("New AppError: %s - %s (%d)", code, message, status)
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

func Wrap(code string, message string, err error, status int) *AppError {
	logger.Error("Wrapped AppError: %s - %s | cause: %v", code, message, err)
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
		Status:  status,
	}
}

func Wrapf(code string, format string, err error, status int, args ...interface{}) *AppError {
	msg := formatMessage(format, args...)
	logger.Error("Wrapped AppError: %s - %s | cause: %v", code, msg, err)
	return &AppError{
		Code:    code,
		Message: msg,
		Err:     err,
		Status:  status,
	}
}

func formatMessage(format string, args ...interface{}) string {
	// Wrap here to avoid fmt directly
	return fmt.Sprintf(format, args...)
}

// ========== Common Presets ==========

func NewInternal(msg string) *AppError {
	return New("internal_error", msg, http.StatusInternalServerError)
}

func NewBadRequest(msg string) *AppError {
	return New("bad_request", msg, http.StatusBadRequest)
}

func NewNotFound(resource string) *AppError {
	return New("not_found", resource+" not found", http.StatusNotFound)
}

func NewUnauthorized() *AppError {
	return New("unauthorized", "unauthorized access", http.StatusUnauthorized)
}

// ========== Helpers ==========

func IsAppError(err error) (*AppError, bool) {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}
