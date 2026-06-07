package roadmap

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

// Get godoc
// @Summary Get roadmap
// @Description Return the IELTS vocabulary roadmap with course, band levels, topics, lessons, and authenticated user progress.
// @Tags Roadmap
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=Response}
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /roadmap [get]
func (h Handler) Get(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}

	result, err := h.service.Get(userID)
	if err != nil {
		writeRoadmapError(c, err)
		return
	}

	response.OK(c, result)
}

func writeRoadmapError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrRoadmapNotFound):
		response.Error(c, http.StatusNotFound, "ROADMAP_NOT_FOUND", "Roadmap was not found")
	case errors.Is(err, ErrUserNotFound):
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
	}
}
