package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"ielts-learning/backend/internal/middleware"
	"ielts-learning/backend/internal/shared/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

func (h Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]string{
			"request": "Invalid JSON request body",
		})
		return
	}

	validatedReq, err := ValidateRegisterRequest(req)
	if err != nil {
		writeAuthError(c, err)
		return
	}

	result, err := h.service.Register(validatedReq)
	if err != nil {
		writeAuthError(c, err)
		return
	}

	response.Created(c, result)
}

func (h Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]string{
			"request": "Invalid JSON request body",
		})
		return
	}

	validatedReq, err := ValidateLoginRequest(req)
	if err != nil {
		writeAuthError(c, err)
		return
	}

	result, err := h.service.Login(validatedReq)
	if err != nil {
		writeAuthError(c, err)
		return
	}

	response.OK(c, result)
}

func (h Handler) Me(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}

	result, err := h.service.Me(userID)
	if err != nil {
		writeAuthError(c, err)
		return
	}

	response.OK(c, result)
}

func writeAuthError(c *gin.Context, err error) {
	var validationErr ValidationError
	switch {
	case errors.As(err, &validationErr):
		response.ValidationError(c, validationErr.Fields)
	case errors.Is(err, ErrEmailAlreadyUsed):
		response.Error(c, http.StatusConflict, "EMAIL_ALREADY_USED", "Email is already registered")
	case errors.Is(err, ErrUsernameAlreadyUsed):
		response.Error(c, http.StatusConflict, "USERNAME_ALREADY_USED", "Username is already registered")
	case errors.Is(err, ErrInvalidCredentials):
		response.Error(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
	case errors.Is(err, ErrUserNotFound):
		response.Error(c, http.StatusNotFound, "USER_NOT_FOUND", "User was not found")
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
	}
}
