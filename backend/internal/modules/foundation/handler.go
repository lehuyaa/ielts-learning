package foundation

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

// Summary godoc
// @Summary Foundation level summary
// @Description Progress through a foundation level: per-skill progress (vocabulary words mastered, grammar mastery, sound accuracy, listening/reading/writing/speaking scores), situation lessons and weekly checkpoints done, streak, the next lesson and — for A1 — the "A1 reached" checklist (spec §1).
// @Tags Foundation
// @Produce json
// @Security BearerAuth
// @Param level query string false "CEFR level (default A1)"
// @Success 200 {object} response.SuccessResponse{data=SummaryResponse}
// @Failure 404 {object} response.ErrorResponse
// @Router /foundation/summary [get]
func (h Handler) Summary(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}
	result, err := h.service.Summary(userID, c.Query("level"))
	if err != nil {
		if errors.Is(err, ErrLevelNotFound) {
			response.Error(c, http.StatusNotFound, "LEVEL_NOT_FOUND", "This level is not in the foundation course yet")
			return
		}
		log.Printf("foundation summary: %v", err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
		return
	}
	response.OK(c, result)
}
