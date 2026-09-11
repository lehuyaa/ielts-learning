package aiconversation

import (
	"errors"
	"log"
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

// CreateScenario godoc
// @Summary Create a conversation scenario
// @Description Create a custom AI Conversation practice scenario owned by the authenticated user.
// @Tags AI Conversation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateScenarioRequest true "Create scenario request"
// @Success 201 {object} response.SuccessResponse{data=ScenarioResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /ai-conversations/scenarios [post]
func (h Handler) CreateScenario(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}

	var req CreateScenarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]string{
			"request": "Invalid JSON request body",
		})
		return
	}

	validatedReq, err := ValidateCreateScenarioRequest(req)
	if err != nil {
		writeScenarioError(c, err)
		return
	}

	result, err := h.service.CreateScenario(userID, validatedReq)
	if err != nil {
		writeScenarioError(c, err)
		return
	}

	response.Created(c, result)
}

// ListScenarios godoc
// @Summary List conversation scenarios
// @Description Return the authenticated user's AI Conversation scenarios.
// @Tags AI Conversation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=ListScenariosResponse}
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /ai-conversations/scenarios [get]
func (h Handler) ListScenarios(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}

	result, err := h.service.ListScenarios(userID)
	if err != nil {
		writeScenarioError(c, err)
		return
	}

	response.OK(c, result)
}

// SendChatMessage godoc
// @Summary Send a chat message to the AI
// @Description Send a message (with optional prior turns for context) and get back the AI's reply.
// @Tags AI Conversation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChatRequest true "Chat request"
// @Success 200 {object} response.SuccessResponse{data=ChatResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 502 {object} response.ErrorResponse
// @Failure 503 {object} response.ErrorResponse
// @Router /ai-conversations/chat [post]
func (h Handler) SendChatMessage(c *gin.Context) {
	_, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}

	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]string{
			"request": "Invalid JSON request body",
		})
		return
	}

	validatedReq, err := ValidateChatRequest(req)
	if err != nil {
		writeScenarioError(c, err)
		return
	}

	result, err := h.service.SendChatMessage(c.Request.Context(), validatedReq)
	if err != nil {
		writeChatError(c, err)
		return
	}

	response.OK(c, result)
}

func writeChatError(c *gin.Context, err error) {
	log.Printf("ai chat error: %v", err)

	if errors.Is(err, ErrAIChatNotConfigured) {
		response.Error(c, http.StatusServiceUnavailable, "AI_CHAT_NOT_CONFIGURED", "AI chat is not configured on the server")
		return
	}

	response.Error(c, http.StatusBadGateway, "AI_CHAT_REQUEST_FAILED", "Unable to get a response from the AI right now")
}

func writeScenarioError(c *gin.Context, err error) {
	var validationErr ValidationError
	switch {
	case errors.As(err, &validationErr):
		response.ValidationError(c, validationErr.Fields)
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
	}
}
