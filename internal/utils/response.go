package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// SuccessResponse returns a success response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse returns an error response
func ErrorResponse(c *gin.Context, statusCode int, message string, err interface{}) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// BadRequestResponse returns a bad request response
func BadRequestResponse(c *gin.Context, message string, err interface{}) {
	ErrorResponse(c, http.StatusBadRequest, message, err)
}

// NotFoundResponse returns a not found response
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message, nil)
}

// UnauthorizedResponse returns an unauthorized response
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message, nil)
}

// ForbiddenResponse returns a forbidden response
func ForbiddenResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, message, nil)
}

// InternalServerErrorResponse returns a internal server error response
func InternalServerErrorResponse(c *gin.Context, err interface{}) {
	ErrorResponse(c, http.StatusInternalServerError, "Internal server error", err)
}

// ValidationErrorResponse returns a validation error response
func ValidationErrorResponse(c *gin.Context, err interface{}) {
	ErrorResponse(c, http.StatusUnprocessableEntity, "Validation error", err)
}

// CreatedResponse returns a created response
func CreatedResponse(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, http.StatusCreated, message, data)
}

// OKResponse returns an OK response
func OKResponse(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, http.StatusOK, message, data)
}

// NoContentResponse returns a no content response
func NoContentResponse(c *gin.Context) {
	c.Status(http.StatusNoContent)
} 