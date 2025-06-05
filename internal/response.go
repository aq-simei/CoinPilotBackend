package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorData  `json:"error,omitempty"`
}

type ErrorData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func JSON(c *gin.Context, statusCode int, success bool, data interface{}, err *ErrorData) {
	c.JSON(statusCode, Response{
		Success: success,
		Data:    data,
		Error:   err,
	})
}

func Success(c *gin.Context, data interface{}) {
	JSON(c, http.StatusOK, true, data, nil)
}

func Created(c *gin.Context, data interface{}) {
	JSON(c, http.StatusCreated, true, data, nil)
}

func BadRequest(c *gin.Context, message string) {
	if message == "" {
		message = "Bad request"
	}
	JSON(c, http.StatusBadRequest, false, nil, &ErrorData{
		Code:    http.StatusBadRequest,
		Message: message,
	})
}

func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized access"
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
		Success: false,
		Error: &ErrorData{
			Code:    http.StatusUnauthorized,
			Message: message,
		},
	})
}

func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "Access forbidden"
	}
	c.AbortWithStatusJSON(http.StatusForbidden, Response{
		Success: false,
		Error: &ErrorData{
			Code:    http.StatusForbidden,
			Message: message,
		},
	})
}

func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	JSON(c, http.StatusNotFound, false, nil, &ErrorData{
		Code:    http.StatusNotFound,
		Message: message,
	})
}

func InternalServerError(c *gin.Context, message string) {
	if message == "" {
		message = "Internal server error"
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
		Success: false,
		Error: &ErrorData{
			Code:    http.StatusInternalServerError,
			Message: message,
		},
	})
}

func CustomError(c *gin.Context, statusCode int, message string) {
	JSON(c, statusCode, false, nil, &ErrorData{
		Code:    statusCode,
		Message: message,
	})
}
