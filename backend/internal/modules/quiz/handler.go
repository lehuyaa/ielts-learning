package quiz

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
// @Summary Get lesson quiz
// @Description Return quiz questions and answer options for a lesson without exposing correct answers.
// @Tags Quiz
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param lessonId path int true "Lesson ID"
// @Success 200 {object} response.SuccessResponse{data=GetResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /lessons/{lessonId}/quiz [get]
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
		writeQuizError(c, err)
		return
	}

	response.OK(c, result)
}

// Submit godoc
// @Summary Submit lesson quiz
// @Description Grade selected answers on the backend, save quiz attempt records, and update authenticated user lesson progress.
// @Tags Quiz
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param lessonId path int true "Lesson ID"
// @Param request body SubmitRequest true "Selected quiz answers"
// @Success 200 {object} response.SuccessResponse{data=SubmitResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /lessons/{lessonId}/quiz/submit [post]
func (h Handler) Submit(c *gin.Context) {
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

	var request SubmitRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.ValidationError(c, map[string]string{
			"answers": "Answers are required.",
		})
		return
	}

	result, err := h.service.Submit(userID, lessonID, request)
	if err != nil {
		writeQuizError(c, err)
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

func writeQuizError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrLessonNotFound):
		response.Error(c, http.StatusNotFound, "LESSON_NOT_FOUND", "Lesson was not found")
	case errors.Is(err, ErrQuizNotFound):
		response.Error(c, http.StatusNotFound, "QUIZ_NOT_FOUND", "Quiz was not found")
	case errors.Is(err, ErrUserNotFound):
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
	case errors.Is(err, ErrLessonLocked):
		response.Error(c, http.StatusForbidden, "LESSON_LOCKED", "Complete previous lessons to unlock this lesson")
	case errors.Is(err, ErrInvalidQuizAnswers):
		response.ValidationError(c, map[string]string{
			"answers": "Answers must reference valid lesson questions and options without duplicates.",
		})
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
	}
}
