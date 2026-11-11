package practice

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

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

// Generate godoc
// @Summary Practice a weakness
// @Description Returns today's practice set for a skill tag (cached per user, tag and day). Grammar / vocabulary tags: 10 AI-generated exercises, each validated like an import (answer key grades as correct, A1 sentences ≤ 12 words). Sound, stress, dictation and number tags: 10 lesson exercises, last answered wrong first. Reading / listening comprehension and speaking are practised in lessons (400).
// @Tags Practice
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body GenerateRequest true "Skill tag, e.g. grammar:g1-to-be"
// @Success 200 {object} response.SuccessResponse{data=SetResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Failure 429 {object} response.ErrorResponse
// @Failure 503 {object} response.ErrorResponse
// @Router /practice/generate [post]
func (h Handler) Generate(c *gin.Context) {
	userID, ok := userIDOf(c)
	if !ok {
		return
	}
	var req GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.SkillTag) == "" {
		response.ValidationError(c, map[string]string{"skillTag": "skillTag is required"})
		return
	}
	result, err := h.service.Generate(c.Request.Context(), userID, req.SkillTag)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, result)
}

// Get godoc
// @Summary Get a practice set
// @Description Exercises (numbered 1..n) and the answers given so far.
// @Tags Practice
// @Produce json
// @Security BearerAuth
// @Param setId path int true "Practice set ID"
// @Success 200 {object} response.SuccessResponse{data=SetResponse}
// @Failure 404 {object} response.ErrorResponse
// @Router /practice/sets/{setId} [get]
func (h Handler) Get(c *gin.Context) {
	userID, setID, ok := params(c)
	if !ok {
		return
	}
	result, err := h.service.Get(userID, setID)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, result)
}

// Answer godoc
// @Summary Answer a practice exercise
// @Description Grades one exercise of the set (exerciseId = its number). Updates grammar mastery and the mistake log; each exercise is answered once.
// @Tags Practice
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param setId path int true "Practice set ID"
// @Param request body AnswerRequest true "Exercise number and answer"
// @Success 200 {object} response.SuccessResponse{data=AnswerResult}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /practice/sets/{setId}/answers [post]
func (h Handler) Answer(c *gin.Context) {
	userID, setID, ok := params(c)
	if !ok {
		return
	}
	var req AnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.ExerciseID == 0 || len(req.Answer) == 0 {
		response.ValidationError(c, map[string]string{"answer": "exerciseId and answer are required"})
		return
	}
	result, err := h.service.Answer(userID, setID, req)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, result)
}

func userIDOf(c *gin.Context) (uint, bool) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
	}
	return userID, ok
}

func params(c *gin.Context) (uint, uint, bool) {
	userID, ok := userIDOf(c)
	if !ok {
		return 0, 0, false
	}
	id, err := strconv.ParseUint(c.Param("setId"), 10, 64)
	if err != nil || id == 0 {
		response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "Invalid practice set id")
		return 0, 0, false
	}
	return userID, uint(id), true
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidTag):
		response.ValidationError(c, map[string]string{"skillTag": "skillTag is required (max 80 characters)"})
	case errors.Is(err, ErrNotPracticable):
		response.Error(c, http.StatusBadRequest, "PRACTICE_NOT_SUPPORTED", "This skill is practised in lessons (reading, listening, speaking)")
	case errors.Is(err, ErrNotEnough):
		response.Error(c, http.StatusUnprocessableEntity, "PRACTICE_NOT_ENOUGH", "Could not build enough exercises for this skill. Try again later.")
	case errors.Is(err, ErrLimitReached):
		response.Error(c, http.StatusTooManyRequests, "AI_LIMIT_REACHED", "You have used today's AI practice. Try again tomorrow.")
	case errors.Is(err, ErrAIUnavailable):
		response.Error(c, http.StatusServiceUnavailable, "AI_UNAVAILABLE", "AI practice is not available right now. Please try again.")
	case errors.Is(err, ErrSetNotFound):
		response.Error(c, http.StatusNotFound, "PRACTICE_NOT_FOUND", "Practice set not found")
	case errors.Is(err, ErrItemNotFound):
		response.Error(c, http.StatusNotFound, "EXERCISE_NOT_FOUND", "Exercise not in this set")
	case errors.Is(err, ErrAlreadyAnswer):
		response.Error(c, http.StatusConflict, "ALREADY_ANSWERED", "This exercise was already answered")
	case errors.Is(err, lessonactivity.ErrInvalidAnswerShape):
		response.Error(c, http.StatusBadRequest, "INVALID_ANSWER_SHAPE", "The answer does not match this exercise type")
	default:
		log.Printf("practice: %v", err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
	}
}
