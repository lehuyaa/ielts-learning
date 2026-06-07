package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Data  any       `json:"data,omitempty"`
	Error *APIError `json:"error,omitempty"`
}

type SuccessResponse struct {
	Data any `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error *APIError `json:"error,omitempty"`
}

type APIError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, APIResponse{Data: data})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, APIResponse{Data: data})
}

func Error(c *gin.Context, status int, code string, message string) {
	c.JSON(status, APIResponse{
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}

func ValidationError(c *gin.Context, fields map[string]string) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Error: &APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Validation failed",
			Fields:  fields,
		},
	})
}
