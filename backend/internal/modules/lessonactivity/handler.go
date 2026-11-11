package lessonactivity

import (
	"errors"
	"log"
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

// ListActivities godoc
// @Summary List lesson activities
// @Description Return a lesson's skill activities in order with the user's progress and lock state. Exercise data is not included.
// @Tags Activities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param lessonId path int true "Lesson ID"
// @Success 200 {object} response.SuccessResponse{data=LessonActivitiesResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /lessons/{lessonId}/activities [get]
func (h Handler) ListActivities(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}
	lessonID, ok := parseIDParam(c, "lessonId")
	if !ok {
		return
	}

	result, err := h.service.ListActivities(userID, lessonID)
	if err != nil {
		writeActivityError(c, err)
		return
	}

	response.OK(c, result)
}

// GetActivity godoc
// @Summary Get activity
// @Description Return an activity with its exercises (prompt and public payload only — never the answers) and the answers already given in the current attempt.
// @Tags Activities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param activityId path int true "Activity ID"
// @Success 200 {object} response.SuccessResponse{data=ActivityDetailResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /activities/{activityId} [get]
func (h Handler) GetActivity(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}
	activityID, ok := parseIDParam(c, "activityId")
	if !ok {
		return
	}

	result, err := h.service.GetActivity(userID, activityID)
	if err != nil {
		writeActivityError(c, err)
		return
	}

	response.OK(c, result)
}

// SubmitAnswer godoc
// @Summary Answer an exercise
// @Description Grade one exercise answer and return immediate feedback. The answer shape depends on the exercise type (docs/foundation-spec.md §5.1).
// @Tags Activities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param activityId path int true "Activity ID"
// @Param request body SubmitAnswerRequest true "Answer"
// @Success 200 {object} response.SuccessResponse{data=AnswerResult}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /activities/{activityId}/answers [post]
func (h Handler) SubmitAnswer(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}
	activityID, ok := parseIDParam(c, "activityId")
	if !ok {
		return
	}

	var req SubmitAnswerRequest
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

	result, err := h.service.SubmitAnswer(c.Request.Context(), userID, activityID, req)
	if err != nil {
		writeActivityError(c, err)
		return
	}

	response.OK(c, result)
}

// CompleteActivity godoc
// @Summary Finish an activity attempt
// @Description Score the current attempt, pass or fail the activity, award XP on the first pass, and complete the lesson when every required activity is passed.
// @Tags Activities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param activityId path int true "Activity ID"
// @Success 200 {object} response.SuccessResponse{data=CompleteActivityResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /activities/{activityId}/complete [post]
func (h Handler) CompleteActivity(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}
	activityID, ok := parseIDParam(c, "activityId")
	if !ok {
		return
	}

	result, err := h.service.CompleteActivity(userID, activityID)
	if err != nil {
		writeActivityError(c, err)
		return
	}

	response.OK(c, result)
}

func parseIDParam(c *gin.Context, name string) (uint, bool) {
	value, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || value == 0 {
		response.ValidationError(c, map[string]string{name: name + " must be a positive integer"})
		return 0, false
	}
	return uint(value), true
}

func writeActivityError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrLessonNotFound):
		response.Error(c, http.StatusNotFound, "LESSON_NOT_FOUND", "Lesson was not found")
	case errors.Is(err, ErrActivityNotFound):
		response.Error(c, http.StatusNotFound, "ACTIVITY_NOT_FOUND", "Activity was not found")
	case errors.Is(err, ErrLessonLocked):
		response.Error(c, http.StatusForbidden, "LESSON_LOCKED", "Complete previous lessons to unlock this lesson")
	case errors.Is(err, ErrActivityLocked):
		response.Error(c, http.StatusForbidden, "ACTIVITY_LOCKED", "Complete the previous activities first")
	case errors.Is(err, ErrActivityNotStarted):
		response.Error(c, http.StatusBadRequest, "ACTIVITY_NOT_STARTED", "Answer at least one exercise before finishing")
	case errors.Is(err, ErrExerciseNotInActivity):
		response.Error(c, http.StatusBadRequest, "EXERCISE_NOT_IN_ACTIVITY", "This exercise does not belong to the activity")
	case errors.Is(err, ErrInvalidAnswerShape):
		response.Error(c, http.StatusBadRequest, "INVALID_ANSWER_SHAPE", "The answer does not match this exercise type")
	case errors.Is(err, ErrNotWritingActivity):
		response.Error(c, http.StatusBadRequest, "NOT_WRITING_ACTIVITY", "This activity is not a writing activity")
	case errors.Is(err, ErrWritingTooShort), errors.Is(err, ErrWritingTooLong):
		response.ValidationError(c, map[string]string{"text": err.Error()})
	case errors.Is(err, ErrWritingNotFound):
		response.Error(c, http.StatusNotFound, "WRITING_NOT_FOUND", "Writing submission was not found")
	case errors.Is(err, ErrWritingNotRewritable):
		response.Error(c, http.StatusConflict, "WRITING_NOT_REWRITABLE", "This submission does not need a rewrite")
	case errors.Is(err, ErrWritingNotCompleted):
		response.Error(c, http.StatusBadRequest, "WRITING_NOT_COMPLETED", "Submit your text and its rewrite before finishing")
	case errors.Is(err, ErrSpeakingNotCompleted):
		response.Error(c, http.StatusBadRequest, "SPEAKING_NOT_COMPLETED", "Finish a conversation with the AI before completing this activity")
	case errors.Is(err, ErrAILimitReached):
		response.Error(c, http.StatusTooManyRequests, "AI_LIMIT_REACHED", "You have used today's AI gradings. Try again tomorrow.")
	case errors.Is(err, ErrAIUnavailable):
		log.Printf("lessonactivity ai: %v", err)
		response.Error(c, http.StatusServiceUnavailable, "AI_UNAVAILABLE", "AI grading is not available right now")
	case errors.Is(err, ErrInvalidAnswerKey):
		log.Printf("lessonactivity: %v", err)
		response.Error(c, http.StatusInternalServerError, "EXERCISE_DATA_INVALID", "This exercise is misconfigured")
	default:
		log.Printf("lessonactivity: %v", err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
	}
}

// SubmitWriting godoc
// @Summary Submit a text for a WRITING activity
// @Description The AI scores the text (task, grammar, vocabulary, coherence; overall = 10 × average) and marks each mistake with its position, correction and a Vietnamese explanation. Mistakes are logged for review. A text with mistakes must be rewritten before the activity can be completed.
// @Tags Activities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SubmitWritingRequest true "Activity and text"
// @Success 200 {object} response.SuccessResponse{data=WritingSubmissionResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 429 {object} response.ErrorResponse
// @Failure 503 {object} response.ErrorResponse
// @Router /writing/submissions [post]
func (h Handler) SubmitWriting(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}
	var req SubmitWritingRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.ActivityID == 0 {
		response.ValidationError(c, map[string]string{"request": "activityId and text are required"})
		return
	}
	result, err := h.service.SubmitWriting(c.Request.Context(), userID, req)
	if err != nil {
		writeActivityError(c, err)
		return
	}
	response.OK(c, result)
}

// RewriteWriting godoc
// @Summary Submit the rewrite of a text
// @Description Checks the corrected text; the submission is COMPLETED when no mistakes remain (otherwise the remaining ones are returned in rewriteReport).
// @Tags Activities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param submissionId path int true "Writing submission ID"
// @Param request body RewriteWritingRequest true "Rewritten text"
// @Success 200 {object} response.SuccessResponse{data=WritingSubmissionResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 429 {object} response.ErrorResponse
// @Router /writing/submissions/{submissionId}/rewrite [post]
func (h Handler) RewriteWriting(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}
	submissionID, ok := parseIDParam(c, "submissionId")
	if !ok {
		return
	}
	var req RewriteWritingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]string{"request": "text is required"})
		return
	}
	result, err := h.service.RewriteWriting(c.Request.Context(), userID, submissionID, req)
	if err != nil {
		writeActivityError(c, err)
		return
	}
	response.OK(c, result)
}
