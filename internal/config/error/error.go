package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/aq-simei/coin-pilot/internal/config/logger"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code int, message string) *AppError {
	logger.Warn("New AppError: %d - %s", code, message)
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func Wrap(code int, message string, err error) *AppError {
	logger.Error("Wrapped AppError: %d - %s | cause: %v", code, message, err)
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func formatMessage(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func NewInternal(msg string) *AppError {
	return New(http.StatusInternalServerError, msg)
}

func NewBadRequest(msg string) *AppError {
	return New(http.StatusBadRequest, msg)
}

func NewNotFound(resource string) *AppError {
	return New(http.StatusNotFound, resource+" not found")
}

func NewUnauthorized() *AppError {
	return New(http.StatusUnauthorized, "unauthorized access")
}

func IsAppError(err error) (*AppError, bool) {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}

