package topic

import (
	"errors"
	"net/http"
	"strconv"

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
// @Summary Get topic detail
// @Description Return topic metadata, parent band level, lesson list, and authenticated user progress.
// @Tags Topics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param topicId path int true "Topic ID"
// @Success 200 {object} response.SuccessResponse{data=Response}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /topics/{topicId} [get]
func (h Handler) Get(c *gin.Context) {
	topicID, err := parseTopicID(c.Param("topicId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "Invalid topic ID")
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}

	result, err := h.service.Get(userID, topicID)
	if err != nil {
		writeTopicError(c, err)
		return
	}

	response.OK(c, result)
}

func parseTopicID(value string) (uint, error) {
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("invalid topic id")
	}

	return uint(id), nil
}

func writeTopicError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrTopicNotFound):
		response.Error(c, http.StatusNotFound, "TOPIC_NOT_FOUND", "Topic was not found")
	case errors.Is(err, ErrUserNotFound):
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
	}
}
