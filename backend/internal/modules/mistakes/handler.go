package mistakes

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

// List godoc
// @Summary Mistake log
// @Description The learner's mistakes from exercises, writing, speaking and pronunciation, grouped by skill tag (most frequent first), with counts and how many review cards are due.
// @Tags Mistakes
// @Produce json
// @Security BearerAuth
// @Param status query string false "open (default: OPEN + REVIEWING) | fixed | all"
// @Success 200 {object} response.SuccessResponse{data=ListResponse}
// @Failure 400 {object} response.ErrorResponse
// @Router /mistakes [get]
func (h Handler) List(c *gin.Context) {
	userID, ok := userIDOf(c)
	if !ok {
		return
	}
	result, err := h.service.List(userID, c.Query("status"))
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, result)
}

// ReviewQueue godoc
// @Summary Due mistake cards
// @Description Mistakes in the review queue (status REVIEWING) that are due now, oldest first.
// @Tags Mistakes
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Max cards (default 20, max 50)"
// @Success 200 {object} response.SuccessResponse{data=ReviewQueueResponse}
// @Router /mistakes/review [get]
func (h Handler) ReviewQueue(c *gin.Context) {
	userID, ok := userIDOf(c)
	if !ok {
		return
	}
	limit, _ := strconv.Atoi(c.Query("limit"))
	result, err := h.service.ReviewQueue(userID, limit)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, result)
}

// Update godoc
// @Summary Change a mistake's status
// @Description REVIEWING adds the mistake to the review queue (a wrong → right card), FIXED marks it learnt, OPEN takes it out of the queue.
// @Tags Mistakes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param mistakeId path int true "Mistake ID"
// @Param request body UpdateRequest true "New status"
// @Success 200 {object} response.SuccessResponse{data=Item}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /mistakes/{mistakeId} [patch]
func (h Handler) Update(c *gin.Context) {
	userID, id, ok := params(c)
	if !ok {
		return
	}
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]string{"status": "status must be OPEN, REVIEWING or FIXED"})
		return
	}
	result, err := h.service.Update(userID, id, req)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, result)
}

// Review godoc
// @Summary Review a mistake card
// @Description Rate a due card (AGAIN, HARD, GOOD, EASY). Three GOOD/EASY reviews in a row mark the mistake FIXED.
// @Tags Mistakes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param mistakeId path int true "Mistake ID"
// @Param request body ReviewRequest true "Rating"
// @Success 200 {object} response.SuccessResponse{data=ReviewResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /mistakes/{mistakeId}/review [post]
func (h Handler) Review(c *gin.Context) {
	userID, id, ok := params(c)
	if !ok {
		return
	}
	var req ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, map[string]string{"rating": "rating must be AGAIN, HARD, GOOD or EASY"})
		return
	}
	result, err := h.service.Review(userID, id, req.Rating)
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
	id, err := strconv.ParseUint(c.Param("mistakeId"), 10, 64)
	if err != nil || id == 0 {
		response.Error(c, http.StatusBadRequest, "INVALID_INPUT", "Invalid mistake id")
		return 0, 0, false
	}
	return userID, uint(id), true
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		response.Error(c, http.StatusNotFound, "MISTAKE_NOT_FOUND", "Mistake not found")
	case errors.Is(err, ErrInvalidStatus):
		response.ValidationError(c, map[string]string{"status": "status must be OPEN, REVIEWING or FIXED (filter: open, fixed, all)"})
	case errors.Is(err, ErrInvalidRating):
		response.ValidationError(c, map[string]string{"rating": "rating must be AGAIN, HARD, GOOD or EASY"})
	case errors.Is(err, ErrNotReviewing):
		response.Error(c, http.StatusConflict, "MISTAKE_NOT_IN_REVIEW", "Add this mistake to review first")
	default:
		log.Printf("mistakes: %v", err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
	}
}
