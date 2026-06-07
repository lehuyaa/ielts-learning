package lesson

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
// @Summary Get lesson detail
// @Description Return lesson metadata, topic and band context, vocabulary list, and authenticated user progress.
// @Tags Lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param lessonId path int true "Lesson ID"
// @Success 200 {object} response.SuccessResponse{data=DetailResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /lessons/{lessonId} [get]
func (h Handler) Get(c *gin.Context) {
	lessonID, err := parseLessonID(c.Param("lessonId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "Invalid lesson ID")
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}

	result, err := h.service.Get(userID, lessonID)
	if err != nil {
		writeLessonError(c, err)
		return
	}

	response.OK(c, result)
}

// Start godoc
// @Summary Start lesson
// @Description Create or update authenticated user lesson progress and mark the lesson as in progress when allowed.
// @Tags Lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param lessonId path int true "Lesson ID"
// @Success 200 {object} response.SuccessResponse{data=StartResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /lessons/{lessonId}/start [post]
func (h Handler) Start(c *gin.Context) {
	lessonID, err := parseLessonID(c.Param("lessonId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "Invalid lesson ID")
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}

	result, err := h.service.Start(userID, lessonID)
	if err != nil {
		writeLessonError(c, err)
		return
	}

	response.OK(c, result)
}

func parseLessonID(value string) (uint, error) {
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("invalid lesson id")
	}

	return uint(id), nil
}

func writeLessonError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrLessonNotFound):
		response.Error(c, http.StatusNotFound, "LESSON_NOT_FOUND", "Lesson was not found")
	case errors.Is(err, ErrUserNotFound):
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
	case errors.Is(err, ErrLessonLocked):
		response.Error(c, http.StatusForbidden, "LESSON_LOCKED", "Complete previous lessons to unlock this lesson")
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
	}
}
