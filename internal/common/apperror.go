package common

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string { return e.Message }

func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func ErrorResponse(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		ApiResponse(c, appErr.Code, false, nil, appErr.Message)
		return
	}
	log.Printf("internal server error: %v", err)
	ApiResponse(c, http.StatusInternalServerError, false, nil, "internal server error")
}
