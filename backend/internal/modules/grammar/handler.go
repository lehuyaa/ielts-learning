package grammar

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"ielts-learning/backend/internal/middleware"
	"ielts-learning/backend/internal/modules/lessonactivity"
	"ielts-learning/backend/internal/shared/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

// List godoc
// @Summary List grammar points
// @Description Return grammar points with the user's mastery (average of the last 20 answers) and the number of practice exercises available.
// @Tags Grammar
// @Produce json
// @Security BearerAuth
// @Param level query string false "CEFR level, e.g. A1"
// @Success 200 {object} response.SuccessResponse{data=ListResponse}
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /grammar-points [get]
func (h Handler) List(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}
	result, err := h.service.List(userID, c.Query("level"))
	if err != nil {
		writeGrammarError(c, err)
		return
	}
	response.OK(c, result)
}

// Get godoc
// @Summary Get grammar point
// @Description Return a grammar theory card with the user's mastery.
// @Tags Grammar
// @Produce json
// @Security BearerAuth
// @Param slug path string true "Grammar point slug"
// @Success 200 {object} response.SuccessResponse{data=DetailResponse}
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /grammar-points/{slug} [get]
func (h Handler) Get(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}
	result, err := h.service.Get(userID, c.Param("slug"))
	if err != nil {
		writeGrammarError(c, err)
		return
	}
	response.OK(c, result)
}

// Practice godoc
// @Summary Get grammar practice
// @Description Return up to 10 exercises for a grammar point from lessons the user has started: recent mistakes first, then new ones. Answers are not included.
// @Tags Grammar
// @Produce json
// @Security BearerAuth
// @Param slug path string true "Grammar point slug"
// @Success 200 {object} response.SuccessResponse{data=PracticeResponse}
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /grammar-points/{slug}/practice [get]
func (h Handler) Practice(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}
	result, err := h.service.Practice(userID, c.Param("slug"))
	if err != nil {
		writeGrammarError(c, err)
		return
	}
	response.OK(c, result)
}

// Answer godoc
// @Summary Answer a grammar practice exercise
// @Description Grade one practice answer, update the user's mastery, and return feedback.
// @Tags Grammar
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param slug path string true "Grammar point slug"
// @Param request body PracticeAnswerRequest true "Answer"
// @Success 200 {object} response.SuccessResponse{data=PracticeAnswerResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /grammar-points/{slug}/practice/answers [post]
func (h Handler) Answer(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}
	var req PracticeAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]string{"request": "Invalid JSON request body"})
		return
	}
	fields := map[string]string{}
	if req.ExerciseID == 0 {
		fields["exerciseId"] = "exerciseId is required"
	}
	if len(req.Answer) == 0 || string(req.Answer) == "null" {
		fields["answer"] = "answer is required"
	}
	if len(fields) > 0 {
		response.ValidationError(c, fields)
		return
	}

	result, err := h.service.Answer(c.Request.Context(), userID, c.Param("slug"), req)
	if err != nil {
		writeGrammarError(c, err)
		return
	}
	response.OK(c, result)
}

func writeGrammarError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrGrammarPointNotFound):
		response.Error(c, http.StatusNotFound, "GRAMMAR_POINT_NOT_FOUND", "Grammar point was not found")
	case errors.Is(err, ErrExerciseNotAvailable):
		response.Error(c, http.StatusBadRequest, "EXERCISE_NOT_AVAILABLE", "This exercise is not available for practice")
	case errors.Is(err, lessonactivity.ErrInvalidAnswerShape):
		response.Error(c, http.StatusBadRequest, "INVALID_ANSWER_SHAPE", "The answer does not match this exercise type")
	default:
		log.Printf("grammar: %v", err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
	}
}
