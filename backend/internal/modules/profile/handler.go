package profile

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

// GetProfile godoc
// @Summary Get current profile
// @Description Return the authenticated user's profile, stats, achievements, and activity summary.
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=GetProfileResponse}
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /me/profile [get]
func (h Handler) GetProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}

	result, err := h.service.GetProfile(userID)
	if err != nil {
		writeProfileError(c, err)
		return
	}

	response.OK(c, result)
}

// UpdateProfile godoc
// @Summary Update current profile
// @Description Update the authenticated user's supported profile fields: name, username, target band, timezone, and locale.
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProfileRequest true "Profile update request"
// @Success 200 {object} response.SuccessResponse{data=UpdateProfileResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /me/profile [patch]
func (h Handler) UpdateProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]string{
			"request": "Invalid JSON request body",
		})
		return
	}

	validatedReq, err := ValidateUpdateProfileRequest(req)
	if err != nil {
		writeProfileError(c, err)
		return
	}

	result, err := h.service.UpdateProfile(userID, validatedReq)
	if err != nil {
		writeProfileError(c, err)
		return
	}

	response.OK(c, result)
}

func writeProfileError(c *gin.Context, err error) {
	var validationErr ValidationError
	switch {
	case errors.As(err, &validationErr):
		response.ValidationError(c, validationErr.Fields)
	case errors.Is(err, ErrUserNotFound):
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
	case errors.Is(err, ErrUsernameAlreadyUsed):
		response.Error(c, http.StatusConflict, "USERNAME_ALREADY_USED", "Username is already registered")
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
	}
}
