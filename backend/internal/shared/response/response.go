package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Data  any       `json:"data,omitempty"`
	Error *APIError `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
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
