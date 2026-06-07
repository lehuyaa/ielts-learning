package flashcard

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

// GetLessonFlashcards godoc
// @Summary Get lesson flashcards
// @Description Return flashcards for a lesson with vocabulary content and authenticated user vocabulary progress.
// @Tags Flashcards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param lessonId path int true "Lesson ID"
// @Success 200 {object} response.SuccessResponse{data=LessonFlashcardsResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /lessons/{lessonId}/flashcards [get]
func (h Handler) GetLessonFlashcards(c *gin.Context) {
	lessonID, err := parseUintParam(c.Param("lessonId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "Invalid lesson ID")
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}

	result, err := h.service.GetLessonFlashcards(userID, lessonID)
	if err != nil {
		writeError(c, err)
		return
	}

	response.OK(c, result)
}

// GetDueReviews godoc
// @Summary Get due reviews
// @Description Return due review flashcards for the authenticated user where next review time is due.
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Maximum number of due cards, max 100"
// @Param topicId query int false "Filter by topic ID"
// @Success 200 {object} response.SuccessResponse{data=DueReviewsResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /reviews/due [get]
func (h Handler) GetDueReviews(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}

	query, err := parseDueReviewQuery(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "Invalid review query")
		return
	}

	result, err := h.service.GetDueReviews(userID, query)
	if err != nil {
		writeError(c, err)
		return
	}

	response.OK(c, result)
}

// Review godoc
// @Summary Save flashcard review
// @Description Save a flashcard rating and update authenticated user vocabulary progress using MVP spaced repetition rules.
// @Tags Flashcards
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ReviewRequest true "Flashcard review request"
// @Success 200 {object} response.SuccessResponse{data=ReviewResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /flashcards/review [post]
func (h Handler) Review(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}

	var request ReviewRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.ValidationError(c, map[string]string{
			"body": "Request body is invalid",
		})
		return
	}

	result, err := h.service.Review(userID, request)
	if err != nil {
		writeError(c, err)
		return
	}

	response.OK(c, result)
}

func parseDueReviewQuery(c *gin.Context) (DueReviewQuery, error) {
	query := DueReviewQuery{Limit: 20}

	if rawLimit := c.Query("limit"); rawLimit != "" {
		limit, err := strconv.Atoi(rawLimit)
		if err != nil {
			return DueReviewQuery{}, err
		}
		query.Limit = limit
	}

	if rawTopicID := c.Query("topicId"); rawTopicID != "" {
		topicID, err := parseUintParam(rawTopicID)
		if err != nil {
			return DueReviewQuery{}, err
		}
		query.TopicID = &topicID
	}

	return query, nil
}

func parseUintParam(value string) (uint, error) {
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("invalid id")
	}

	return uint(id), nil
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrLessonNotFound):
		response.Error(c, http.StatusNotFound, "LESSON_NOT_FOUND", "Lesson was not found")
	case errors.Is(err, ErrVocabularyNotFound):
		response.Error(c, http.StatusNotFound, "VOCABULARY_NOT_FOUND", "Vocabulary was not found")
	case errors.Is(err, ErrUserNotFound):
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
	case errors.Is(err, ErrLessonLocked):
		response.Error(c, http.StatusForbidden, "LESSON_LOCKED", "Complete previous lessons to unlock this lesson")
	case errors.Is(err, ErrInvalidRating):
		response.ValidationError(c, map[string]string{
			"rating": "Rating must be AGAIN, HARD, GOOD, or EASY",
		})
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
	}
}
