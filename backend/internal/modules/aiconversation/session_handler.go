package aiconversation

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"ielts-learning/backend/internal/middleware"
	"ielts-learning/backend/internal/shared/response"
)

const maxMessageChars = 1000

// StartSession godoc
// @Summary Start (or resume) a role-play session
// @Description Start a session for a scenario (system or own). With activityId it counts for that lesson's SPEAKING activity. An unfinished session of the same scenario is resumed unless fresh=true, which abandons it.
// @Tags AI Conversation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body StartSessionRequest true "Scenario slug and optional activity"
// @Success 200 {object} response.SuccessResponse{data=SessionResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /ai-conversations/sessions [post]
func (h Handler) StartSession(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}
	var req StartSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.ScenarioSlug) == "" {
		response.ValidationError(c, map[string]string{"scenarioSlug": "scenarioSlug is required"})
		return
	}
	result, err := h.service.StartSession(userID, req)
	if err != nil {
		writeSessionError(c, err)
		return
	}
	response.OK(c, result)
}

// GetSession godoc
// @Summary Get a session
// @Description Session with its scenario brief, messages and (when finished) the feedback report.
// @Tags AI Conversation
// @Produce json
// @Security BearerAuth
// @Param sessionId path int true "Session ID"
// @Success 200 {object} response.SuccessResponse{data=SessionResponse}
// @Failure 404 {object} response.ErrorResponse
// @Router /ai-conversations/sessions/{sessionId} [get]
func (h Handler) GetSession(c *gin.Context) {
	userID, sessionID, ok := sessionParams(c)
	if !ok {
		return
	}
	result, err := h.service.GetSession(userID, sessionID)
	if err != nil {
		writeSessionError(c, err)
		return
	}
	response.OK(c, result)
}

// SendSessionMessage godoc
// @Summary Send one turn in a session
// @Description The AI replies in role with the level-aware prompt (spec §6.1), seeing the last 20 messages. Rejected once every turn is used (finish the session then).
// @Tags AI Conversation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param sessionId path int true "Session ID"
// @Param request body SessionMessageRequest true "Learner message"
// @Success 200 {object} response.SuccessResponse{data=SessionTurnResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 429 {object} response.ErrorResponse
// @Failure 503 {object} response.ErrorResponse
// @Router /ai-conversations/sessions/{sessionId}/messages [post]
func (h Handler) SendSessionMessage(c *gin.Context) {
	userID, sessionID, ok := sessionParams(c)
	if !ok {
		return
	}
	var req SessionMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Message) == "" || len(req.Message) > maxMessageChars {
		response.ValidationError(c, map[string]string{"message": "message is required (max 1000 characters)"})
		return
	}
	result, err := h.service.SendSessionMessage(c.Request.Context(), userID, sessionID, strings.TrimSpace(req.Message))
	if err != nil {
		writeSessionError(c, err)
		return
	}
	response.OK(c, result)
}

// Hint godoc
// @Summary Get a hint
// @Description A Vietnamese hint and one English sentence to say next (spec §6.5). Not allowed in checkpoints.
// @Tags AI Conversation
// @Produce json
// @Security BearerAuth
// @Param sessionId path int true "Session ID"
// @Success 200 {object} response.SuccessResponse{data=HintResponse}
// @Failure 403 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /ai-conversations/sessions/{sessionId}/hint [post]
func (h Handler) Hint(c *gin.Context) {
	userID, sessionID, ok := sessionParams(c)
	if !ok {
		return
	}
	result, err := h.service.Hint(c.Request.Context(), userID, sessionID)
	if err != nil {
		writeSessionError(c, err)
		return
	}
	response.OK(c, result)
}

// Finish godoc
// @Summary Finish a session
// @Description The AI writes the feedback report (spec §6.2): 4 scores, overall, mistakes (logged), good phrases, new words, tip. XP when the task was completed.
// @Tags AI Conversation
// @Produce json
// @Security BearerAuth
// @Param sessionId path int true "Session ID"
// @Success 200 {object} response.SuccessResponse{data=SessionResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /ai-conversations/sessions/{sessionId}/finish [post]
func (h Handler) Finish(c *gin.Context) {
	userID, sessionID, ok := sessionParams(c)
	if !ok {
		return
	}
	result, err := h.service.Finish(c.Request.Context(), userID, sessionID)
	if err != nil {
		writeSessionError(c, err)
		return
	}
	response.OK(c, result)
}

func sessionParams(c *gin.Context) (uint, uint, bool) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return 0, 0, false
	}
	id, err := strconv.ParseUint(c.Param("sessionId"), 10, 64)
	if err != nil || id == 0 {
		response.ValidationError(c, map[string]string{"sessionId": "sessionId must be a positive number"})
		return 0, 0, false
	}
	return userID, uint(id), true
}

func writeSessionError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrScenarioNotFound):
		response.Error(c, http.StatusNotFound, "SCENARIO_NOT_FOUND", "Scenario was not found")
	case errors.Is(err, ErrSessionNotFound):
		response.Error(c, http.StatusNotFound, "SESSION_NOT_FOUND", "Session was not found")
	case errors.Is(err, ErrActivityMismatch):
		response.Error(c, http.StatusBadRequest, "ACTIVITY_MISMATCH", "This activity does not use this scenario")
	case errors.Is(err, ErrSessionFinished):
		response.Error(c, http.StatusConflict, "SESSION_FINISHED", "This session is already finished")
	case errors.Is(err, ErrTurnLimit):
		response.Error(c, http.StatusConflict, "SESSION_TURN_LIMIT", "You have used every turn; finish the session to see your report")
	case errors.Is(err, ErrSessionEmpty):
		response.Error(c, http.StatusBadRequest, "SESSION_EMPTY", "Say something before finishing")
	case errors.Is(err, ErrHintsNotAllowed):
		response.Error(c, http.StatusForbidden, "HINTS_NOT_ALLOWED", "Hints are not allowed in this conversation")
	case errors.Is(err, ErrAILimitReached):
		response.Error(c, http.StatusTooManyRequests, "AI_LIMIT_REACHED", "You have used today's AI conversation turns. Try again tomorrow.")
	case errors.Is(err, ErrAIChatNotConfigured):
		response.Error(c, http.StatusServiceUnavailable, "AI_UNAVAILABLE", "AI is not configured (OPENAI_API_KEY)")
	case errors.Is(err, ErrAIUnavailable):
		log.Printf("aiconversation: %v", err)
		response.Error(c, http.StatusServiceUnavailable, "AI_UNAVAILABLE", "The AI did not respond. Please try again.")
	default:
		log.Printf("aiconversation: %v", err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
	}
}
