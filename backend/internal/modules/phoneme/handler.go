package phoneme

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

// List godoc
// @Summary List sound cards
// @Description Return the IPA chart: every phoneme in chart order with its first example word and the user's accuracy.
// @Tags Pronunciation
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=ListResponse}
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /phonemes [get]
func (h Handler) List(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}
	result, err := h.service.List(userID)
	if err != nil {
		writePhonemeError(c, err)
		return
	}
	response.OK(c, result)
}

// Get godoc
// @Summary Get sound card
// @Description Return one sound card (mouth position, example words, minimal pairs, common mistakes), the lessons that teach it and its neighbours on the chart.
// @Tags Pronunciation
// @Produce json
// @Security BearerAuth
// @Param symbol path string true "IPA symbol, URL-encoded (e.g. i%CB%90 for iː)"
// @Success 200 {object} response.SuccessResponse{data=DetailResponse}
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /phonemes/{symbol} [get]
func (h Handler) Get(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication is required")
		return
	}
	result, err := h.service.Get(userID, c.Param("symbol"))
	if err != nil {
		writePhonemeError(c, err)
		return
	}
	response.OK(c, result)
}

func writePhonemeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrPhonemeNotFound):
		response.Error(c, http.StatusNotFound, "PHONEME_NOT_FOUND", "Sound was not found")
	default:
		log.Printf("phoneme: %v", err)
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")
	}
}
